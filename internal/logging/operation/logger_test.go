package operation

import (
	"bytes"
	"errors"
	stdlog "log"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

func TestStarted_LogsOperationLifecycle(t *testing.T) {
	output := NewSUT(t)

	completed := Started("exchange_rate", "service", "get_exchange_rate")
	completed()

	assert.Contains(t, output.String(), "exchange_rate service: operation started operation=get_exchange_rate")
	assert.Contains(t, output.String(), "exchange_rate service: operation completed operation=get_exchange_rate")
	assert.Contains(t, output.String(), "duration=")
}

func TestFailed_LogsDependencyAndError(t *testing.T) {
	output := NewSUT(t)

	Failed("exchange_rate", "service", "get_exchange_rate", "redis", errors.New("connection refused"))

	assert.Contains(t, output.String(), "operation=get_exchange_rate")
	assert.Contains(t, output.String(), "dependency=redis")
	assert.Contains(t, output.String(), `error="connection refused"`)
}

func TestValidationAndNilInput_LogReasons(t *testing.T) {
	output := NewSUT(t)

	ValidationFailed("exchange_rate", "gRPC handler", "get_exchange_rate", "request is nil")
	NilInput("exchange_rate", "mapper", "ranking_to_proto")

	assert.Contains(t, output.String(), "validation failed")
	assert.Contains(t, output.String(), "reason=request is nil")
	assert.Contains(t, output.String(), "nil input")
}

func TestResult_LogsCount(t *testing.T) {
	output := NewSUT(t)

	Result("exchange_rate", "cache", "get_all_exchange_rates", 3)

	assert.Contains(t, output.String(), "operation=get_all_exchange_rates")
	assert.Contains(t, output.String(), "count=3")
}

func NewSUT(t *testing.T) *bytes.Buffer {
	t.Helper()

	previousOutput := stdlog.Writer()
	previousFlags := stdlog.Flags()
	previousPrefix := stdlog.Prefix()

	output := &bytes.Buffer{}
	stdlog.SetOutput(output)
	stdlog.SetFlags(0)
	stdlog.SetPrefix("")

	t.Cleanup(func() {
		stdlog.SetOutput(previousOutput)
		stdlog.SetFlags(previousFlags)
		stdlog.SetPrefix(previousPrefix)
	})

	return output
}
