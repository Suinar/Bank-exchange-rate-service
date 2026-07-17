# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/bank-exchange-rate-service ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app

COPY --from=builder /out/bank-exchange-rate-service /app/bank-exchange-rate-service

EXPOSE 50053

HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/bank-exchange-rate-service", "healthcheck"]

USER nonroot:nonroot

ENTRYPOINT ["/app/bank-exchange-rate-service"]
