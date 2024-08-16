ARG GO_VERSION=1.21
FROM golang:${GO_VERSION} AS build
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server ./cmd/dnevnik-rg
COPY ./config/stage/config.yml /bin/stage/config.yml
COPY ./config/prod/config.yml /bin/prod/config.yml
FROM alpine:latest AS final
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
COPY --from=build /bin/server /bin/
COPY --from=build /bin/stage/config.yml /bin/stage/config.yml
COPY --from=build /bin/prod/config.yml /bin/prod/config.yml
EXPOSE 8000
ENTRYPOINT [ "/bin/server" ]
