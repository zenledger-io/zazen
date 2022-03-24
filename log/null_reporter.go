package log

import "context"

type nullReporter struct{}

func NewNoneReporter() Reporter {
	return nullReporter{}
}

func (r nullReporter) Start(context.Context) error {
	return nil
}

func (r nullReporter) Errorf(string, ...interface{}) error {
	return nil
}

func (r nullReporter) Monitor() {}
