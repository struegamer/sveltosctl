package logging

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

type KlogLogger struct {
	logger     logr.Logger
	ctx        context.Context
	infoLevel  int
	debugLevel int
	warnLevel  int
	verbose    bool
}

func NewKlogLogger(ctx context.Context) *KlogLogger {
	// Init Klog right from start, no need to do this in the main code
	// just initialize this Logger instance.
	klog.InitFlags(nil)
	if ctx == nil {
		ctx = context.TODO()
	}
	return &KlogLogger{
		logger:     klog.Background(),
		ctx:        ctx,
		infoLevel:  0,
		warnLevel:  1,
		debugLevel: 4,
		verbose:    false,
	}
}

func (l *KlogLogger) Logger() logr.Logger {
	return l.logger
}

func (l *KlogLogger) Info(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.V(l.infoLevel).Info(msg, keysAndValues...)
	}

}
func (l *KlogLogger) Debug(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.V(l.debugLevel).Info(msg, keysAndValues...)
	}

}
func (l *KlogLogger) Warn(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.V(l.warnLevel).Info(msg, keysAndValues...)
	}
}
func (l *KlogLogger) Error(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.Error(nil, msg, keysAndValues...)
	}

}
func (l *KlogLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.Error(msg, keysAndValues...)
}

func (l *KlogLogger) SetVerbose(toggle bool) {
	l.verbose = toggle
}

//func (l *KlogLogger) makeListFromMap(fields map[string]interface{}) []interface{} {
//	var result = make([]interface{}, 0)
//	for k, v := range fields {
//		result = append(result, k, v)
//	}
//	return result
//}
