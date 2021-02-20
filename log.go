package main

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func initLogger() {
	logger = logrus.New()
	if debugMode {
		logger.SetLevel(logrus.DebugLevel)
	}
}
