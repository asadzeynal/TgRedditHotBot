package util

type Logger interface {
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
}
