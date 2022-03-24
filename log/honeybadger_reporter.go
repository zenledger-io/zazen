package log

import (
	"context"
	"fmt"

	"github.com/honeybadger-io/honeybadger-go"
)

const (
	errMessageTag = "_error_message"
	levelTag      = "severity"
)

var (
	HoneybadgerAPIKey string
	HoneybadgerSync   bool
)

type honeybadgerReporter struct{}

func (r honeybadgerReporter) Start(ctx context.Context) error {
	cfg := getHoneybadgerConfig(ctx)
	if cfg.APIKey == "" {
		return fmt.Errorf("cannot start honeybadger without an api key")
	}

	honeybadger.Configure(cfg)
	honeybadger.BeforeNotify(func(notice *honeybadger.Notice) error {
		var errMessage string
		for i, tag := range notice.Tags {
			if tag == errMessageTag && len(notice.Tags) > i+1 {
				errMessage = notice.Tags[i+1]
				notice.Tags = append(notice.Tags[:i], notice.Tags[i+2:]...)
				break
			}
		}
		if errMessage != "" {
			notice.ErrorMessage = errMessage
		}

		return nil
	})

	return nil
}

func (r honeybadgerReporter) Errorf(format string, args ...interface{}) error {
	return notifyHoneybadger(LevelError, format, args...)
}

func (r honeybadgerReporter) Monitor() {
	honeybadger.Monitor()
}

func notifyHoneybadger(level, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	var notifyErr error

	err := extractErrFromArgs(args)
	if err != nil {
		_, notifyErr = honeybadger.Notify(err, honeybadger.Tags{errMessageTag, msg, levelTag, level})
	} else {
		_, notifyErr = honeybadger.Notify(msg, honeybadger.Tags{levelTag, level})
	}

	return notifyErr
}

func getHoneybadgerConfig(ctx context.Context) honeybadger.Configuration {
	return honeybadger.Configuration{
		Logger: ContextLogger(ctx),
		APIKey: HoneybadgerAPIKey,
		Sync:   HoneybadgerSync,
	}
}

func extractErrFromArgs(args []interface{}) error {
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			return err
		}
	}

	return nil
}
