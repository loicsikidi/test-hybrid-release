FROM cgr.dev/chainguard/static:latest

COPY awesomecli /awesomecli

ENTRYPOINT ["/awesomecli"]
