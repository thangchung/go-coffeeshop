package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	LogError(err error)
	Fatal(message string, args ...interface{})
}

type Logger struct {
	logger *logrus.Entry
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warm":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.NewEntry(logrus.StandardLogger())
	logger.Logger.SetLevel(l)

	return &Logger{logger: logger}
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

func (l *Logger) LogError(err error) {
	l.logger.Error(err)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalf(message, args...)
	os.Exit(1)
}
