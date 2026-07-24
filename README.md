# Bank Exchange Rate Service

A Go service that retrieves Monobank exchange rates, stores them in Redis, and exposes them through gRPC.

## How the Service Works

1. On startup, the service connects to Redis and Kafka.
2. The Monobank client retrieves the current exchange-rate list.
3. The complete snapshot is written to Redis atomically.
4. The service layer reads exchange rates only through `IExchangeRateCache`.
5. The data is available through gRPC methods that:
   - retrieve a single rate by `currencyIsoFrom` and `currencyIsoTo`;
   - retrieve all rates by `currencyIsoFrom`.
6. The snapshot is refreshed automatically at the interval specified by `MONOBANK_REFRESH_INTERVAL`.

The complete Redis snapshot is replaced using an atomic Lua script. Concurrent updates cannot mix different exchange-rate sets, and stale Redis indexes are removed when they are read.

## Main Components

- Go 1.25.5
- gRPC
- Monobank currency API
- Redis 7
- Kafka
- Docker Compose
- Testcontainers for Go
- GoMock and Testify

## Project Structure

```text
cmd/app                                  service entry point
internal/configs                         configuration
internal/delivery/grpc                   gRPC server and handlers
internal/external/monobank               Monobank client
internal/repositories/cahce              Redis connection and initialization
internal/repositories/cahce/exchange_rate Redis exchange-rate cache
internal/services/exchange_rate          service for reading rates from the cache
internal/brokers/kafka                   Kafka topic initialization
internal/mocks                           generated GoMock mocks
internal/test/fixture                    test fixtures
pkg/core                                 domain structures
docker                                   Docker Compose, Dockerfile, and Kafka configuration
```

The `cahce` directory name is retained to match the current project structure.

## Requirements

For running the service with Docker:

- Docker Desktop or Docker Engine;
- Docker Compose.

For local development:

- Go 1.25.5 or a compatible newer version;
- running Redis and Kafka instances;
- GNU Make (optional).

Docker is also required for Redis integration tests. Testcontainers automatically creates and removes an isolated `redis:7-alpine` container.

## Configuration

Create a local `.env` file from the example:

```bash
cp .env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Main environment variables:

| Variable | Default value | Purpose |
|---|---:|---|
| `GRPC_PORT` | `50053` | gRPC server port |
| `GRPC_NETWORK` | `tcp` | gRPC listener network |
| `MONOBANK_BASE_URL` | `https://api.monobank.ua` | Monobank API URL |
| `MONOBANK_CURRENCY_ENDPOINT` | `/bank/currency` | Exchange-rate endpoint |
| `MONOBANK_REFRESH_INTERVAL` | `5m` | Background exchange-rate refresh interval |
| `KAFKA_BROKERS` | none | List of Kafka brokers |
| `GET_RELATIVE_RANKING_REQUEST_TOPIC` | `exchange-rate.get-relative-ranking.request` | Single-rate request topic |
| `GET_ALL_RANKING_REQUEST_TOPIC` | `exchange-rate.get-all-ranking.request` | Exchange-rate list request topic |
| `REDIS_ADDR` | `localhost:6380` | Redis address for local execution |
| `REDIS_PASSWORD` | empty | Redis password |
| `REDIS_DB` | `0` | Redis database number |
| `REDIS_POOL_SIZE` | `10` | Maximum Redis connection pool size |
| `REDIS_MIN_IDLE_CONNS` | `2` | Minimum number of idle connections |

`MONOBANK_REFRESH_INTERVAL` must be greater than zero and supports the Go duration format, such as `30s`, `5m`, or `1h`.

## Running with Docker

Build the images and start all services:

```bash
docker compose -f docker/docker-compose.yml up -d --build
```

Alternatively, use Make:

```bash
make docker-up-build
```

Check the service status:

```bash
make docker-ps
make docker-health
```

View logs:

```bash
make docker-app-logs
make docker-redis-logs
make docker-kafka-logs
```

Stop the environment:

```bash
make docker-down
```

Remove the containers together with the Redis and Kafka volumes:

```bash
make docker-reset
```

## Running Locally

Start Redis and Kafka first, then specify their addresses in `.env`.

```bash
go run ./cmd/app
```

Alternatively:

```bash
make run
```

During startup, the application:

- verifies the Redis connection;
- retrieves the initial snapshot from Monobank;
- creates the required Kafka topics;
- starts the background exchange-rate refresh process;
- starts the gRPC server.

If the initial connection or snapshot retrieval fails, the application does not start the gRPC server.

## Tests

Run all tests:

```bash
go test -count=1 ./...
```

Alternatively:

```bash
make test
```

Run individual test suites:

```bash
make test-client
make test-service
make test-handler
```

Run additional checks:

```bash
go vet ./...
make test-cover
make test-race
```

The Redis integration tests are located in:

```text
internal/repositories/cahce/exchange_rate
```

They use Testcontainers for Go, start a real Redis instance on a random port, and automatically remove the container after completion.

## Redis Data Model

Each currency pair is stored as a Redis hash. Additional Redis sets are used as indexes:

- by the complete pair key, `currencyIsoFrom:currencyIsoTo`;
- by the source currency, `currencyIsoFrom`;
- to track all keys in the current snapshot.

Supported operations:

```go
AddAllExchangeRate(ctx, rates)
GetExchangeRate(ctx, currencyIsoFrom, currencyIsoTo)
GetAllExchangeRate(ctx, currencyIsoFrom)
```

`AddAllExchangeRate` completely replaces the previous snapshot, including when an empty list is provided.

## Health Check

The same binary supports a health-check mode:

```bash
./bank-exchange-rate-service healthcheck
```

The health check connects to the local gRPC server and verifies the standard gRPC health status.
