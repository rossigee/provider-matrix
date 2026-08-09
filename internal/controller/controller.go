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

// Package controller contains the controllers for the Matrix provider.
package controller

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/powerlevel"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/providerconfig"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/room"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/roomalias"
	"github.com/crossplane-contrib/provider-matrix/internal/controller/user"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Setup sets up all controllers for the Matrix provider.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	if err := providerconfig.Setup(mgr); err != nil {
		return err
	}
	if err := user.Setup(mgr, o); err != nil {
		return err
	}
	if err := room.Setup(mgr, o); err != nil {
		return err
	}
	if err := powerlevel.Setup(mgr, o); err != nil {
		return err
	}
	if err := roomalias.Setup(mgr, o); err != nil {
		return err
	}
	return nil
}