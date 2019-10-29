package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

// WithLogger creates a new context with provided logger.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger returns the logger from the current context. Creates a new one if not present
func GetLogger(ctx context.Context) *logrus.Entry {
	var logger *logrus.Entry

	// Get a logger, if it is present.
	loggerInterface := ctx.Value(loggerKey{})
	if loggerInterface != nil {
		if lgr, ok := loggerInterface.(*logrus.Entry); ok {
			logger = lgr
		}
	}

	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return logger
}
