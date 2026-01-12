package logging

import (
	"context"

	"github.com/go-logr/logr"
)

type CliLogger struct {
	logger logr.Logger
}

func NewLogger() *CliLogger {
	return &CliLogger{
		logger: InitLogger(),
	}
}

func NewLoggerFromContext(ctx context.Context) *CliLogger {
	return &CliLogger{
		logger: LoggerFromContext(ctx),
	}
}

func (l *CliLogger) Logger() logr.Logger {
	return l.logger
}
