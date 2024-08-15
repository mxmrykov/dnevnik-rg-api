package logger

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
)

func NewLogger() *zerolog.Logger {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	return &logger
}
