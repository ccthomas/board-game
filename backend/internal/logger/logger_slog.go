package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// https://github.com/golang/go/blob/master/src/log/slog/example_custom_levels_test.go
// Exported constants from a custom logging package.
const (
	LevelTrace     slog.Level = slog.Level(-8)
	LevelDebug                = slog.LevelDebug
	LevelInfo                 = slog.LevelInfo
	LevelNotice               = slog.Level(2) // TODO Implement in interface
	LevelWarning              = slog.LevelWarn
	LevelError                = slog.LevelError
	LevelEmergency            = slog.Level(12) // TODO Implement in interface
)

type LoggerSlog struct {
	logger *slog.Logger
	file   *os.File
}

func NewLoggerSlog() (*LoggerSlog, error) {
	logLevel := LevelError
	levelStr := os.Getenv("LOG_LEVEL")
	switch levelStr {
	case "TRACE":
		logLevel = LevelTrace
	case "NOTICE":
		logLevel = LevelNotice
	case "EMERGENCY":
		logLevel = LevelEmergency
	default:
		var parsed slog.Level
		err := parsed.UnmarshalText([]byte(levelStr))
		if err != nil {
			return nil, err
		}
		logLevel = parsed
	}

	logFileName := os.Getenv("LOG_FILE_NAME")
	if logFileName == "" {
		logFileName = "logfile.log"
	}

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	var writer io.Writer = file

	logToTerminal := os.Getenv("LOG_TO_TERMINAL")
	if logToTerminal == "true" {
		writer = io.MultiWriter(file, os.Stdout)
	}

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		// Set a custom level to show all log output. The default value is
		// LevelInfo, which would drop Debug and Trace logs.
		Level: logLevel,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			// Customize the name of the level key and the output string, including
			// custom level values.
			if a.Key == slog.LevelKey {
				// Rename the level key from "level" to "sev".
				a.Key = "sev"

				// Handle custom level values.
				level := a.Value.Any().(slog.Level)

				// TODO Apply recommendation from guide: https://github.com/golang/go/blob/master/src/log/slog/example_custom_levels_test.go
				// This could also look up the name from a map or other structure, but
				// this demonstrates using a switch statement to rename levels. For
				// maximum performance, the string values should be constants, but this
				// example uses the raw strings for readability.
				switch {
				case level < LevelDebug:
					a.Value = slog.StringValue("TRACE")
				case level < LevelInfo:
					a.Value = slog.StringValue("DEBUG")
				case level < LevelNotice:
					a.Value = slog.StringValue("INFO")
				case level < LevelWarning:
					a.Value = slog.StringValue("NOTICE")
				case level < LevelError:
					a.Value = slog.StringValue("WARNING")
				case level < LevelEmergency:
					a.Value = slog.StringValue("ERROR")
				default:
					a.Value = slog.StringValue("EMERGENCY")
				}
			}

			return a
		},
	})

	baseLogger := slog.New(handler)
	logger := &LoggerSlog{
		logger: baseLogger,
		file:   file,
	}

	logger.WithFields("file_name", "logger_slog.go").Info(
		"Configured logger with settings",
		"log_level", levelStr,
		"log_file_name", logFileName,
		"log_terminal", logToTerminal,
	)

	return logger, nil
}

func (l *LoggerSlog) WithFields(args ...any) Logger {
	return &LoggerSlog{
		logger: l.logger.With(args...),
		file:   l.file,
	}
}

func (l *LoggerSlog) Trace(msg string, args ...any) {
	ctx := context.Background()
	l.logger.Log(ctx, LevelTrace, msg, args...)
}

func (l *LoggerSlog) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *LoggerSlog) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *LoggerSlog) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *LoggerSlog) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
