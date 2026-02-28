// Package logger...TODO
package logger

type Fields map[string]interface{}

type Logger interface {
	WithFields(fields ...any) Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Trace(msg string, args ...any)
	Warn(msg string, args ...any)
}
