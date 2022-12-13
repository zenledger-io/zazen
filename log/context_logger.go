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

func (l contextLogger) Debugf(format string, args ...any) {
	l.logToFunc(Debugf, format, args)
}

func (l contextLogger) Printf(format string, args ...any) {
	l.logToFunc(Printf, format, args)
}

func (l contextLogger) Warnf(format string, args ...any) {
	l.logToFunc(Warnf, format, args)
}

func (l contextLogger) Errorf(format string, args ...any) {
	l.logToFunc(Errorf, format, args)
}

func (l contextLogger) DebugT(msg string, tags ...Tag) {
	l.logToFuncT(DebugT, msg, tags)
}

func (l contextLogger) PrintT(msg string, tags ...Tag) {
	l.logToFuncT(PrintT, msg, tags)
}

func (l contextLogger) WarnT(msg string, tags ...Tag) {
	l.logToFuncT(WarnT, msg, tags)
}

func (l contextLogger) ErrorT(msg string, tags ...Tag) {
	l.logToFuncT(ErrorT, msg, tags)
}

func (l contextLogger) Monitor() {
	defer reporter.Monitor()

	if err := recover(); err != nil {
		l.Printf("%v", err)
		panic(err)
	}
}

func (l contextLogger) logToFunc(f func(string, ...any), format string, args []any) {
	format, args = l.formatArgs(format, args)
	f(format, args...)
}

func (l contextLogger) logToFuncT(f func(string, ...Tag), msg string, tags []Tag) {
	tags = append(tags, NewTag("context_uuid", l.uuid))
	f(msg, tags...)
}

func (l contextLogger) formatArgs(format string, args []any) (string, []any) {
	if l.uuid == "" {
		return format, args
	}

	newArgs := make([]any, len(args)+1)
	newArgs[0] = l.uuid
	for i, arg := range args {
		newArgs[i+1] = arg
	}

	return fmt.Sprintf("%s %s", "[%s]", format), newArgs
}
