package log

import (
	"context"
	"log"
)

const (
	DefaultFlags = log.LUTC | log.Ldate | log.Ltime | log.Lmsgprefix

	LevelError = "error"
)

var (
	reporter = NewNoneReporter()
)

func init() {
	log.SetFlags(DefaultFlags)
}

func Start(ctx context.Context) {
	var hbReporter honeybadgerReporter
	if err := hbReporter.Start(ctx); err != nil {
		ContextLogger(ctx).Printf("%v", err)
	} else {
		reporter = hbReporter
	}
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	if err := reporter.Errorf(format, args...); err != nil {
		log.Printf("reporting error: %v", err)
	}

	log.Printf(format, args...)
}
