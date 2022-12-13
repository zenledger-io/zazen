package log

type Logger interface {
	Debugf(string, ...any)
	Printf(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)

	DebugT(string, ...Tag)
	PrintT(string, ...Tag)
	WarnT(string, ...Tag)
	ErrorT(string, ...Tag)

	Monitor()
}
