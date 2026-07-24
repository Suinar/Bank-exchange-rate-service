package config

import stdlog "log"

func Loading(file string) {
	stdlog.Printf("config: loading file=%s", file)
}
func FileLoaded(file string) {
	stdlog.Printf("config: file loaded file=%s", file)
}
func FileNotFound(file string) {
	stdlog.Printf("config: file not found, using environment variables file=%s", file)
}
func LoadingFailed(file string, err error) {
	stdlog.Printf("config: loading failed file=%s error=%q", file, err)
}
func Loaded(grpcPort, redisAddress, monobankURL string) {
	stdlog.Printf(
		"config: loaded grpc_port=%s redis_address=%s monobank_url=%s",
		grpcPort, redisAddress, monobankURL,
	)
}
