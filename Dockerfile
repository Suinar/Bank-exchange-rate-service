# ---------- Builder ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o bank-exchange-rate-service \
    ./cmd/server

# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bank-exchange-rate-service .

USER nonroot:nonroot

EXPOSE 50053

ENTRYPOINT ["/bank-exchange-rate-service"]