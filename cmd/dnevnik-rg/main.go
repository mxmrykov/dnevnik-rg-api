package main

import (
	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/app"
	"log"
)

func main() {
	appConfig, errNewConfig := config.NewConfig()
	if errNewConfig != nil {
		log.Fatalf("cannot start app: config error; %v", errNewConfig)
		return
	}
	app.App(appConfig)
}
