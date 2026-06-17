package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	// человекочитаемый вывод в консоль (в проде можно убрать ConsoleWriter — будет чистый JSON)
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	Log = zerolog.New(output).With().Timestamp().Logger()
}
