package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

func ContextLogger(ctx context.Context) Logger {
	return contextLogger{uuid: ContextUUID(ctx)}
}

type contextLogger struct {
	uuid string
}

func (l contextLogger) Printf(format string, args ...interface{}) {
	l.logToFunc(zapLogger.Info, format, args)
}

func (l contextLogger) Errorf(format string, args ...interface{}) {
	l.logToFunc(zapLogger.Error, format, args)
}

func (l contextLogger) Monitor() {
	defer reporter.Monitor()

	if err := recover(); err != nil {
		l.Printf("%v", err)
		panic(err)
	}
}

func (l contextLogger) logToFunc(f func(string, ...zap.Field), format string, args []interface{}) {
	format, args = l.formatArgs(format, args)
	f(fmt.Sprintf(format, args...))
}

func (l contextLogger) formatArgs(format string, args []interface{}) (string, []interface{}) {
	if l.uuid == "" {
		return format, args
	}

	newArgs := make([]interface{}, len(args)+1)
	newArgs[0] = l.uuid
	for i, arg := range args {
		newArgs[i+1] = arg
	}

	return fmt.Sprintf("%s %s", "[%s]", format), newArgs
}
