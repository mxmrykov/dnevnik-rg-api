package main

import (
	"context"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/app"
	"dnevnik-rg.ru/pkg/clients/vault"
	zerolog "dnevnik-rg.ru/pkg/logger"
	"dnevnik-rg.ru/pkg/utils"

	log "github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	vaultCfg, err := config.NewVaultConfig()

	if err != nil {
		log.Err(err).Send()
		return
	}

	vaultClient, err := vault.NewVaultClient(vaultCfg)

	if err != nil {
		log.Err(err).Send()
		return
	}

	if err = utils.ExtractVaultDataToENV(ctx, vaultClient, vaultCfg); err != nil {
		log.Err(err).Send()
		return
	}

	appConfig, err := config.NewConfig()
	if err != nil {
		log.Err(err).Send()
		return
	}

	logger := zerolog.NewLogger()

	logger.Info().Msg("starting service...")

	app.App(appConfig, logger)
}
