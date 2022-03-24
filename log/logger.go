package log

type Logger interface {
	Printf(string, ...interface{})
	Errorf(string, ...interface{})
	Monitor()
}
