/*
Copyright 2025 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	"github.com/crossplane-contrib/provider-matrix/apis"
	"github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/powerlevel"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/room"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/roomalias"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/user"
	"github.com/crossplane-contrib/provider-matrix/internal/features"
	"github.com/crossplane-contrib/provider-matrix/internal/version"
)

func main() {
	var (
		app                        = kingpin.New(filepath.Base(os.Args[0]), "Matrix support for Crossplane.").DefaultEnvars()
		debug                      = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncInterval               = app.Flag("sync", "Sync interval controls how often all resources will be double-checked for drift.").Default("1h").Duration()
		pollInterval               = app.Flag("poll", "Poll interval controls how often an individual resource should be checked for drift.").Default("1m").Duration()
		maxReconcileRate           = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may checked for drift from the desired state.").Default("100").Int()
		leaderElection             = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		namespace                  = app.Flag("namespace", "Namespace used to set as default scope in default secret store config.").Default("crossplane-system").Envar("POD_NAMESPACE").String()
		enableExternalSecretStores = app.Flag("enable-external-secret-stores", "Enable support for ExternalSecretStores.").Default("false").Envar("ENABLE_EXTERNAL_SECRET_STORES").Bool()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-matrix"))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	log.Info("Provider starting up",
		"provider", "provider-matrix",
		"version", version.Version,
		"go-version", runtime.Version(),
		"platform", runtime.GOOS+"/"+runtime.GOARCH,
		"sync-interval", syncInterval.String(),
		"poll-interval", pollInterval.String(),
		"max-reconcile-rate", *maxReconcileRate,
		"leader-election", *leaderElection,
		"namespace", *namespace,
		"external-secret-stores", *enableExternalSecretStores,
		"debug-mode", *debug)

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	// Feature flags
	o := controller.Options{
		Logger:                  log,
		MaxConcurrentReconciles: *maxReconcileRate,
		PollInterval:            *pollInterval,
		GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
		Features:                &feature.Flags{},
	}
	if *enableExternalSecretStores {
		o.Features.Enable(features.EnableAlphaExternalSecretStores)
		log.Info("Alpha feature enabled", "flag", features.EnableAlphaExternalSecretStores)
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:             *leaderElection,
		LeaderElectionID:           "crossplane-leader-election-provider-matrix",
		Cache:                      cache.Options{SyncPeriod: syncInterval},
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add Matrix APIs to scheme")

	ctx := context.Background()

	// Initialize default ProviderConfig if it doesn't exist
	if err := createDefaultProviderConfig(ctx, mgr, *namespace); err != nil {
		log.Debug("Cannot create default ProviderConfig", "error", err)
	}

	kingpin.FatalIfError(user.Setup(mgr, o), "Cannot setup User controller")
	kingpin.FatalIfError(room.Setup(mgr, o), "Cannot setup Room controller")
	kingpin.FatalIfError(powerlevel.Setup(mgr, o), "Cannot setup PowerLevel controller")
	kingpin.FatalIfError(roomalias.Setup(mgr, o), "Cannot setup RoomAlias controller")

	kingpin.FatalIfError(mgr.AddHealthzCheck("healthz", healthz.Ping), "Cannot add health check")
	kingpin.FatalIfError(mgr.AddReadyzCheck("readyz", healthz.Ping), "Cannot add ready check")

	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

func createDefaultProviderConfig(ctx context.Context, mgr ctrl.Manager, namespace string) error {
	pc := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{
							Name:      "matrix-creds",
							Namespace: namespace,
						},
						Key: "credentials",
					},
				},
			},
			HomeserverURL: "https://matrix.org",
		},
	}

	err := mgr.GetClient().Create(ctx, pc)
	if kerrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
