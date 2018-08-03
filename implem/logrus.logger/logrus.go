package logger

import (
	"log"

	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	env    string
	Logger *logrus.Logger
}

type CredentialsGetter interface {
	GetCredentials() string
}

func NewLogger(env, logLevel, logFormat string) uc.Logger {
	logger := logrus.New()
	l, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithField("level", logLevel).Warn("Invalid log level, fallback to 'info'")
	} else {
		logrus.SetLevel(l)
	}

	switch logFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	return &LogrusLogger{env: env, Logger: logger}
}

func (l LogrusLogger) Log(args ...interface{}) {
	if l.Logger == nil {
		return
	}

	// for stackDriver the log should contain 2 args :
	// first : the error (of uc.Error type),
	// second : an additional message (most likely the request that lead to the error)
	if len(args) == 2 {
		castedError, ok := args[0].(error)
		if ok {
			l.newLog(castedError, args[1])
		} else {
			l.Logger.Info(args...)
		}
		return
	}
	l.Logger.Info(args...)
}

//fixme
func (l LogrusLogger) newLog(err error, usecase interface{}) {
	l.Logger.WithError(err).Error(usecase)
}

type SimpleLogger struct{}

func (l SimpleLogger) Log(args ...interface{}) {
	log.Println(args...)
}
