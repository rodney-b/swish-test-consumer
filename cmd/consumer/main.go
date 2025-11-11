package main

import (
	"log/slog"
	"os"

	"github.com/rodney-b/swish-test-consumer/internal/app/consumer"
	"github.com/rodney-b/swish-test-consumer/internal/pkg/config"
	"github.com/rodney-b/swish-test-consumer/internal/pkg/logger"
)

func main() {
	appConfig, err := config.InitAppConfig()
	if err != nil {
		errLog := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("package", "main")
		errLog.Error("error initializing the application config", "error", err.Error())
		return
	}

	logger.Initialize(appConfig)
	log := logger.New("main")

	err = consumer.Run(appConfig)
	if err != nil {
		log.Error("consumer error", "error", err.Error())
		return
	}
}
