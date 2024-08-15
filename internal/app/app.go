package app

import (
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/store"
	"dnevnik-rg.ru/pkg/http"
	"dnevnik-rg.ru/pkg/postgres"
	"github.com/rs/zerolog"
)

func App(appConfig *config.Config, logger *zerolog.Logger) {
	postgresConnection, errConnectPostgres := postgres.NewPostgres(&appConfig.Postgres)
	if errConnectPostgres != nil {
		logger.Err(errConnectPostgres).Msg("cannot connect to postrges")
		return
	}
	rgStore := store.NewStore(postgresConnection, 20*time.Second)

	logger.Info().Msg("database connected")

	http.NewHttp(&appConfig.Http, rgStore, true, logger)
}
