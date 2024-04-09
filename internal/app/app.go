package app

import (
	"log"
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/store"
	"dnevnik-rg.ru/pkg/http"
	"dnevnik-rg.ru/pkg/postgres"
)

func App(appConfig *config.Config) {
	postgresConnection, errConnectPostgres := postgres.NewPostgres(&appConfig.Postgres)
	if errConnectPostgres != nil {
		log.Fatalf("cannot connect to postrges: %v", errConnectPostgres)
		return
	}
	rgStore := store.NewStore(postgresConnection, 20*time.Second)
	http.NewHttp(&appConfig.Http, rgStore, true)
}
