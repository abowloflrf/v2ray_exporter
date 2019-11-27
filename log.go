package main

import (
	"log"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func initLogger() {
	devLoggerConfig := zap.NewDevelopmentConfig()
	devLoggerConfig.DisableStacktrace = true
	if !debugMode {
		devLoggerConfig.Level.SetLevel(zap.InfoLevel)
	}
	logger, err := devLoggerConfig.Build()
	if err != nil {
		log.Fatal("init logger", err)
	}
	sugar = logger.Sugar()
}
