package utils

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

// NewLogger creates the main context and injects logger into it.
func NewLogger() (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = ctxlogrus.ToContext(ctx, logrus.NewEntry(logrus.StandardLogger()))

	return ctx, cancelFunc
}
