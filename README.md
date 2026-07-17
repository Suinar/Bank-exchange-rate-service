# Bank Exchange Rate Service

Go microservice that retrieves current exchange rates from the public Monobank API and exposes them to the banking platform over gRPC.

## Technologies

- Go 1.25
- gRPC and Protocol Buffers
- Monobank HTTP API
- Docker and Docker Compose
- Apache Kafka in KRaft mode (optional local infrastructure; integration is not implemented yet)
- Viper
- GoMock
- Testify

## Features

- Fetches the current Monobank currency snapshot
- Returns one exchange rate by currency pair
- Returns all available rates for a source currency
- Propagates request cancellation to the upstream HTTP call
- Uses bounded upstream and healthcheck timeouts
- Exposes the standard gRPC health service
- Supports graceful container shutdown
- Uses a non-root distroless runtime image
- Keeps successful and error tests in separate files

## Project structure

```text
.
├── cmd/app/                                # Application entry point and healthcheck
├── internal/
│   ├── configs/                            # Environment and optional .env configuration
│   ├── delivery/grpc/                      # gRPC server
│   │   └── handler/exchange_rate/          # Exchange-rate transport handlers and tests
│   ├── external/monobank/                  # Monobank HTTP client and transport models
│   ├── mocks/                              # Generated GoMock implementations
│   ├── services/exchange_rate/             # Application logic and tests
│   └── test/                               # Shared constants and fixtures
├── pkg/core/                               # Domain models
├── docker-compose.yml                      # Application and optional Kafka profiles
├── Dockerfile                              # Multi-stage non-root image
├── .env.example                            # Configuration template
└── Makefile                                # Development commands
```

Test files follow the same convention as `Bank-repository-service`:

- `*_test.go` contains successful behavior;
- `*_error_test.go` contains validation, dependency, transport, and nil-input errors;
- both files stay beside the implementation in the same package directory.

## Configuration

Copy the local template when environment overrides are needed:

```sh
cp .env.example .env
```

The `.env` file is optional and ignored by Git. Environment variables take precedence.

| Variable | Default | Purpose |
|---|---:|---|
| `GRPC_PORT` | `50053` | gRPC listen port |
| `GRPC_NETWORK` | `tcp` | Listener network |
| `MONOBANK_BASE_URL` | `https://api.monobank.ua` | Monobank API base URL |
| `MONOBANK_CURRENCY_ENDPOINT` | `/bank/currency` | Currency endpoint path |
| `EXCHANGE_RATE_IMAGE` | `bank-exchange-rate-service:local` | Compose application image |
| `KAFKA_IMAGE` | `bitnami/kafka:latest` | Optional Kafka image |
| `KAFKA_PORT` | `29092` | Kafka host port |
| `KAFKA_ADVERTISED_HOST` | `localhost` | Kafka external advertised host |
| `KAFKA_UI_IMAGE` | `provectuslabs/kafka-ui:latest` | Optional Kafka UI image |
| `KAFKA_UI_PORT` | `8081` | Kafka UI host port |

Kafka settings are infrastructure-only until producer or consumer logic is added.

## Quick start

### Run locally

```sh
go run ./cmd/app
```

The gRPC server listens on `localhost:50053` by default.

### Run in Docker

```sh
docker compose up -d --build exchange-rate-service
docker compose ps
```

Kafka is optional, so a normal application start does not consume resources for an unused broker. Start it only when needed:

```sh
docker compose --profile kafka up -d kafka
```

Start Kafka UI together with its broker:

```sh
docker compose --profile tools up -d kafka-ui
```

## Make commands

Run `make help` for the complete list.

| Command | Description |
|---|---|
| `make run` | Run the service locally |
| `make build` | Build the application |
| `make test` | Run all tests |
| `make test-race` | Run tests with the race detector |
| `make test-cover` | Print statement coverage |
| `make docker-up-build` | Build and start the application |
| `make docker-kafka-up` | Start optional Kafka |
| `make docker-tools-up` | Start optional Kafka UI and its dependency |
| `make docker-config` | Validate and render Compose configuration |

## Testing

Run the complete quality check:

```sh
go vet ./...
go test -race -count=1 ./...
```

Run individual layers:

```sh
make test-client
make test-service
make test-handler
```

Generate coverage:

```sh
go test -coverprofile coverage.out ./...
go tool cover -func coverage.out
```

All current tests are self-contained: the Monobank client uses `httptest`, while service and handler dependencies use GoMock. No live Monobank or Kafka connection is required.

## Test coverage

The current service and Monobank client statement-coverage baseline is 100%. Generate and inspect the results independently:

```sh
go test -coverprofile service.coverage ./internal/services/exchange_rate
go tool cover -func service.coverage

go test -coverprofile client.coverage ./internal/external/monobank
go tool cover -func client.coverage
```

### Exchange-rate service

| Function | Coverage |
|---|---:|
| `NewExchangeRateService` | 100.0% |
| `GetExchangeRate` | 100.0% |
| `GetAllExchangeRate` | 100.0% |
| `MapperToProto` | 100.0% |
| **Total service coverage** | **100.0%** |

### Monobank client

| Function | Coverage |
|---|---:|
| `NewMonobankClient` | 100.0% |
| `GetAllExchangeRate` | 100.0% |
| `MonobankMapper` | 100.0% |
| `mapMonobankRate` | 100.0% |
| **Total client coverage** | **100.0%** |

Generated mocks, shared fixtures, bootstrap packages, and the application entry point are excluded from these layer-level figures because their package percentages do not represent service or client behavior.

## gRPC API

The service implements `RankingRepository` from `github.com/Suinar/Bank-proto`:

- `GetExchangeRate` returns a rate for one source/target currency pair;
- `GetAllExchangeRate` returns all rates for a source currency.

Protocol definitions are versioned through the Go module dependency.

## Architecture

```text
gRPC client
    │
    ▼
exchange-rate handler
    │
    ▼
exchange-rate service
    │
    ▼
Monobank HTTP client ──► Monobank API
```

The handler validates transport input, the service filters and maps domain data, and the external client owns HTTP behavior and upstream response mapping.

## Docker image and health

Build the image manually:

```sh
docker build -t bank-exchange-rate-service:local .
```

The final image runs as a non-root user, has a read-only filesystem in Compose, and contains only the statically linked application binary. The same binary provides a bounded gRPC healthcheck:

```sh
/app/bank-exchange-rate-service healthcheck
```

The container reports `unhealthy` when the local gRPC health service is unreachable or not serving. On `SIGTERM`, readiness changes to `NOT_SERVING`, active calls receive time to finish, and the server is force-stopped if graceful termination cannot complete.
