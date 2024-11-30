FROM golang:1.23.3-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/server ./cmd/server/

FROM debian:12.8-slim  AS debian-golang-dev

RUN adduser nonroot

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    apt-get update && apt-get install -y ca-certificates

COPY --chown=nonroot --from=builder /app /app

USER nonroot:nonroot

WORKDIR /app

EXPOSE 8080

CMD ["/app/server"]
