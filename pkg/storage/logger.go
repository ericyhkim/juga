package storage

type Logger interface {
	Warn(format string, args ...any)
}

type nopLogger struct{}

func (n *nopLogger) Warn(format string, args ...any) {}

func NewNopLogger() Logger {
	return &nopLogger{}
}
