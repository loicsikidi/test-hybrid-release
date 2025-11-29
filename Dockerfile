FROM cgr.dev/chainguard/static:latest

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/awesomecli /awesomecli

ENTRYPOINT ["/awesomecli"]
