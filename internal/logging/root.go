package logging

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

func InitLogger() logr.Logger {
	klog.InitFlags(nil)
	return klog.NewKlogr()
}

func InitLoggerWithContext(ctx context.Context) context.Context {
	klog.InitFlags(nil)
	ctx = klog.NewContext(ctx, klog.Background())
	return ctx
}

func LoggerFromContext(ctx context.Context) logr.Logger {
	logger := klog.FromContext(ctx)
	return logger
}
