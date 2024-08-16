package utils

import (
	"context"
	"os"
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/pkg/clients/vault"
)

func ExtractVaultDataToENV(ctx context.Context, client *vault.VaultClient, cfg *config.VaultCfg) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	appJwtSecret, err := client.GetVaultData(ctx, cfg.VaultSecret.App.Path, cfg.VaultSecret.App.JwtSecretVariable)
	if err != nil {
		return err
	}
	if err = os.Setenv("APP_JWT_SECRET", appJwtSecret); err != nil {
		return err
	}

	postgres_username, err := client.GetVaultData(ctx, cfg.VaultSecret.PostgresVault.Path, cfg.VaultSecret.PostgresVault.UsernameVariable)
	if err != nil {
		return err
	}
	if err = os.Setenv("PG_USER", postgres_username); err != nil {
		return err
	}

	postgres_password, err := client.GetVaultData(ctx, cfg.VaultSecret.PostgresVault.Path, cfg.VaultSecret.PostgresVault.PasswordVariable)
	if err != nil {
		return err
	}
	if err = os.Setenv("PG_PASSWORD", postgres_password); err != nil {
		return err
	}

	telebot_token, err := client.GetVaultData(ctx, cfg.VaultSecret.Telebot.Path, cfg.VaultSecret.Telebot.TokenVariable)
	if err != nil {
		return err
	}
	if err = os.Setenv("TELEBOT_TOKEN", telebot_token); err != nil {
		return err
	}

	return nil
}
