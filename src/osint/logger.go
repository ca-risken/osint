package main

import (
	"github.com/sirupsen/logrus"
)

var (
	appLogger = newAppLogger()
)

func newAppLogger() *logrus.Logger {
	appLogger := logrus.New()
	appLogger.Formatter = &logrus.JSONFormatter{}
	return appLogger
}
