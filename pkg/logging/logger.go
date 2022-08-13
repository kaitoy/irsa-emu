package logging

import (
	"github.com/sirupsen/logrus"
	kwlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
)

var baseLogger kwlog.Logger

func init() {
	baseLogger = kwlog.Noop
}

// Init initialize the base logger.
func Init(level logrus.Level) {
	logrusLog := logrus.New()
	logrusLogEntry := logrus.NewEntry(logrusLog).WithField("app", "irsa-emu-webhook")
	logrusLogEntry.Logger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogEntry.Logger.SetLevel(level)
	baseLogger = kwlogrus.NewLogrus(logrusLogEntry)
}

// GetLogger returns the base logger.
func GetLogger() kwlog.Logger {
	return baseLogger
}
