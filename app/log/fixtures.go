package log

import (
	"strings"
	"sync"
	"testing"

	zerolog "github.com/rs/zerolog/log"
)

var mut sync.Mutex

func CaptureLogging(t *testing.T, callback func()) string {
	mut.Lock()

	t.Helper()

	SetupLogger()
	origLogger := zerolog.Logger
	defer func() {
		zerolog.Logger = zerolog.Logger.Output(origLogger)
		mut.Unlock()
	}()

	buf := strings.Builder{}
	zerolog.Logger = zerolog.Logger.Output(&buf)
	callback()

	return buf.String()
}
