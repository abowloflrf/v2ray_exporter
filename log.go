package main

import (
	"log"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	devLoggerConfig := zap.NewDevelopmentConfig()
	devLoggerConfig.DisableStacktrace = true

	logger, err := devLoggerConfig.Build()
	if err != nil {
		log.Fatal("init logger", err)
	}
	sugar = logger.Sugar()
}
