package couchbase

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/couchbase/gocbcore.v7"
)

type Logger struct {
	Level gocbcore.LogLevel
}

func NewLogger(level string) *Logger {
	var l gocbcore.LogLevel

	switch level {
	case "error":
		l = gocbcore.LogError
	case "warn":
		l = gocbcore.LogWarn
	case "info":
		l = gocbcore.LogInfo
	case "debug":
		l = gocbcore.LogDebug
	case "trace":
		l = gocbcore.LogTrace
	case "sched":
		l = gocbcore.LogSched
	case "max":
		l = gocbcore.LogMaxVerbosity
	}

	return &Logger{l}
}

func (l *Logger) Log(level gocbcore.LogLevel, offset int, format string, v ...interface{}) error {
	if level > l.Level {
		return nil
	}

	log := logrus.WithFields(logrus.Fields{
		"module":          "gocb",
		"gocbcore.offset": offset,
	})

	switch level {
	case gocbcore.LogError:
		log.Errorf(format, v...)
	case gocbcore.LogWarn:
		log.Warnf(format, v...)
	case gocbcore.LogInfo:
		log.Infof(format, v...)
	case gocbcore.LogDebug:
		log.WithField("gocbcore.LogLevel", "Debug").Debugf(format, v...)
	case gocbcore.LogTrace:
		log.WithField("gocbcore.LogLevel", "Trace").Debugf(format, v...)
	case gocbcore.LogSched:
		log.WithField("gocbcore.LogLevel", "Sched").Debugf(format, v...)
	case gocbcore.LogMaxVerbosity:
		log.WithField("gocbcore.LogLevel", "MaxVerbosity").Debugf(format, v...)
	}

	return nil
}
