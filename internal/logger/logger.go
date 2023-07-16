package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	Log = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
}
