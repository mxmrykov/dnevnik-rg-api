package main

import (
	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/app"
	zerolog "dnevnik-rg.ru/pkg/logger"

	"log"
)

func main() {
	appConfig, errNewConfig := config.NewConfig()
	if errNewConfig != nil {
		log.Fatalf("cannot start app: config error; %v", errNewConfig)
		return
	}

	logger := zerolog.NewLogger()

	logger.Info().Msg("starting service...")

	app.App(appConfig, logger)
}
