package logging

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2/textlogger"
)

type KlogTextLogger struct {
	logger     logr.Logger
	ctx        context.Context
	infoLevel  int
	debugLevel int
	warnLevel  int
	verbose    bool
}

func NewKlogTextLogger(ctx context.Context) *KlogTextLogger {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &KlogTextLogger{
		logger:     textlogger.NewLogger(textlogger.NewConfig()),
		ctx:        ctx,
		infoLevel:  0,
		warnLevel:  1,
		debugLevel: 4,
	}
}

func (l *KlogTextLogger) Logger() logr.Logger {
	return l.logger
}

func (l *KlogTextLogger) Info(msg string, keysAndValues ...interface{}) {

	if l.verbose {
		l.logger.V(l.infoLevel).Info(msg, keysAndValues...)
	}
}
func (l *KlogTextLogger) Debug(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.V(l.debugLevel).Info(msg, keysAndValues...)
	}
}
func (l *KlogTextLogger) Warn(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.V(l.warnLevel).Info(msg, keysAndValues...)
	}
}
func (l *KlogTextLogger) Error(msg string, keysAndValues ...interface{}) {
	if l.verbose {
		l.logger.Error(nil, msg, keysAndValues...)
	}
}
func (l *KlogTextLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.Error(msg, keysAndValues...)
}
func (l *KlogTextLogger) SetVerbose(toggle bool) {
	l.verbose = toggle
}

//func (l *KlogLogger) makeListFromMap(fields map[string]interface{}) []interface{} {
//	var result = make([]interface{}, 0)
//	for k, v := range fields {
//		result = append(result, k, v)
//	}
//	return result
//}
