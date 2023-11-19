package log

import (
	"strings"
	"testing"

	zerolog "github.com/rs/zerolog/log"
)

func CaptureLogging(t *testing.T, callback func()) string {
	t.Helper()

	SetupLogger()
	origLogger := zerolog.Logger
	defer func() {
		zerolog.Logger = zerolog.Logger.Output(origLogger)
	}()

	buf := strings.Builder{}
	zerolog.Logger = zerolog.Logger.Output(&buf)
	callback()

	return buf.String()
}
