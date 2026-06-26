#syntax=docker/dockerfile:1
FROM golang:1.26.4-trixie@sha256:76a29248dedcd75870e95cbd90cc8cb356db082404ac7d3a5803f276c3ba79c9 AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod

RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=${GOMODCACHE} \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=${GOMODCACHE} \
    --mount=type=bind,target=. \
    go build -o /usr/bin/app main.go

FROM gcr.io/distroless/static-debian13:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240

WORKDIR /app

COPY --from=builder /usr/bin/app /usr/bin/app

CMD [ "/usr/bin/app" ]
