package app

import (
	stdlog "log"
	time "time"
)

func Configure() {
	stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds | stdlog.LUTC)
}

func Started()              { stdlog.Printf("application: starting") }
func Stopped()              { stdlog.Printf("application: stopped") }
func ConfigurationLoading() { stdlog.Printf("application: loading configuration") }
func ConfigurationFailed(err error) {
	stdlog.Printf("application: configuration failed error=%q", err)
}
func ConfigurationLoaded(grpcAddress, redisAddress string, refreshInterval time.Duration) {
	stdlog.Printf(
		"application: configuration loaded grpc_address=%s redis_address=%s refresh_interval=%s",
		grpcAddress, redisAddress, refreshInterval,
	)
}
func RedisConnecting(address string) {
	stdlog.Printf("application: connecting to Redis address=%s", address)
}
func RedisConnectionFailed(address string, err error) {
	stdlog.Printf("application: Redis connection failed address=%s error=%q", address, err)
}
func RedisConnected(address string) {
	stdlog.Printf("application: Redis connected address=%s", address)
}
func ComponentsInitializing() {
	stdlog.Printf("application: initializing cache, services, and gRPC handlers")
}
func ComponentsInitialized() {
	stdlog.Printf("application: cache, services, and gRPC handlers initialized")
}
func GRPCServerStarting(address string) {
	stdlog.Printf("application: gRPC server starting address=%s", address)
}
func GRPCServerFailed(err error) {
	stdlog.Printf("application: gRPC server failed error=%q", err)
}
func GRPCServerStopped() { stdlog.Printf("application: gRPC server stopped") }
func ShutdownSignalReceived() {
	stdlog.Printf("application: shutdown signal received")
}
func InitialExchangeRatesLoading() {
	stdlog.Printf("application: loading initial exchange rates")
}
func InitialExchangeRatesLoadingFailed(err error) {
	stdlog.Printf("application: initial exchange-rate loading failed error=%q", err)
}
func InitialExchangeRatesLoaded() {
	stdlog.Printf("application: initial exchange rates loaded")
}
func RefreshSchedulerStarted(interval time.Duration) {
	stdlog.Printf("application: exchange-rate refresh scheduler started interval=%s", interval)
}
func RefreshSchedulerStopped() {
	stdlog.Printf("application: exchange-rate refresh scheduler stopped")
}
func RefreshFailed(duration time.Duration, err error) {
	stdlog.Printf("application: exchange-rate refresh failed duration=%s error=%q", duration, err)
}
func RefreshCompleted(duration time.Duration) {
	stdlog.Printf("application: exchange-rate refresh completed duration=%s", duration)
}
