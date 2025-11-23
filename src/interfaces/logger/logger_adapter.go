// Package logger
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"metar-provider/src/interfaces/config"
)

type Decorator struct {
	logger           Interface
	loggerPrefixName string
}

func NewLoggerAdapter(
	logger Interface,
	loggerPrefixName string,
) *Decorator {
	return &Decorator{
		logger:           logger,
		loggerPrefixName: loggerPrefixName,
	}
}

func (loggerDecorator *Decorator) Init(logPath, logName, logLevel string, logConfig *config.LogConfig) {
	loggerDecorator.logger.Init(logPath, logName, logLevel, logConfig)
}

func (loggerDecorator *Decorator) ShutdownCallback(ctx context.Context) error {
	return loggerDecorator.logger.ShutdownCallback(ctx)
}

func (loggerDecorator *Decorator) LogHandler() *slog.Logger {
	return loggerDecorator.logger.LogHandler()
}

func (loggerDecorator *Decorator) Debug(msg string) {
	loggerDecorator.logger.Debug(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *Decorator) Debugf(msg string, v ...interface{}) {
	loggerDecorator.Debug(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *Decorator) Info(msg string) {
	loggerDecorator.logger.Info(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *Decorator) Infof(msg string, v ...interface{}) {
	loggerDecorator.Info(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *Decorator) Warn(msg string) {
	loggerDecorator.logger.Warn(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *Decorator) Warnf(msg string, v ...interface{}) {
	loggerDecorator.Warn(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *Decorator) Error(msg string) {
	loggerDecorator.logger.Error(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *Decorator) Errorf(msg string, v ...interface{}) {
	loggerDecorator.Error(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *Decorator) Fatal(msg string) {
	loggerDecorator.logger.Fatal(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *Decorator) Fatalf(msg string, v ...interface{}) {
	loggerDecorator.Debug(fmt.Sprintf(msg, v...))
}
