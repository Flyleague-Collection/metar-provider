// Package logger
package logger

import (
	"context"
	"log/slog"
	"metar-provider/src/interfaces/config"
)

type Interface interface {
	Init(logPath, logName, logLevel string, logConfig *config.LogConfig)
	ShutdownCallback(ctx context.Context) error
	LogHandler() *slog.Logger
	Debug(msg string)
	Debugf(msg string, v ...interface{})
	Info(msg string)
	Infof(msg string, v ...interface{})
	Warn(msg string)
	Warnf(msg string, v ...interface{})
	Error(msg string)
	Errorf(msg string, v ...interface{})
	Fatal(msg string)
	Fatalf(msg string, v ...interface{})
}
