FROM gcr.io/distroless/static@sha256:3d0f463de06b7ddff27684ec3bfd0b54a425149d0f8685308b1fdf297b0265e9

# OpenContainers Image Spec Labels
LABEL org.opencontainers.image.title="provider-matrix"
LABEL org.opencontainers.image.description="Crossplane provider for Matrix homeserver management"
LABEL org.opencontainers.image.vendor="rossigee"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.source="https://github.com/rossigee/provider-matrix"
LABEL org.opencontainers.image.url="https://github.com/rossigee/provider-matrix"
LABEL org.opencontainers.image.documentation="https://github.com/rossigee/provider-matrix/blob/master/README.md"

ARG TARGETOS
ARG TARGETARCH

ADD bin/$TARGETOS\_$TARGETARCH/provider /usr/local/bin/provider

USER 65532
ENTRYPOINT ["provider"]