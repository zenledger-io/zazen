package log

import (
	"context"
	"fmt"
)

func ContextLogger(ctx context.Context) Logger {
	return contextLogger{uuid: ContextUUID(ctx)}
}

type contextLogger struct {
	uuid string
}

func (l contextLogger) Printf(format string, args ...interface{}) {
	l.logToFunc(Printf, format, args)
}

func (l contextLogger) Errorf(format string, args ...interface{}) {
	l.logToFunc(Errorf, format, args)
}

func (l contextLogger) Monitor() {
	defer reporter.Monitor()

	if err := recover(); err != nil {
		l.Printf("%v", err)
		panic(err)
	}
}

func (l contextLogger) logToFunc(f func(string, ...interface{}), format string, args []interface{}) {
	format, args = l.formatArgs(format, args)
	f(format, args...)
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
