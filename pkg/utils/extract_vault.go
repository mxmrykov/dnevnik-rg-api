package utils

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/pkg/clients/vault"
)

var ENV_NAMES = map[string]string{
	"UsernameVariable":  "PG_USER",
	"PasswordVariable":  "PG_PASS",
	"JwtSecretVariable": "APP_JWT_SECRET",
	"TokenVariable":     "TELEBOT_TOKEN",
}

func ExtractVaultDataToENV(ctx context.Context, cl *vault.VaultClient, cfg *config.VaultCfg) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c := reflect.ValueOf(cfg.VaultSecret)

	for i := range c.NumField() {
		if c.Field(i).Kind() == reflect.Struct {
			innerVal := c.Field(i)

			tp := innerVal.Type()
			for j := range innerVal.NumField() {
				if j == 0 {
					continue
				}

				v, err := cl.GetVaultData(ctx, innerVal.Field(0).Interface().(string), innerVal.Field(j).Interface().(string))
				if err != nil {
					return err
				}

				if err = os.Setenv(ENV_NAMES[tp.Field(j).Name], v); err != nil {
					return err
				}
			}
		}
	}

	//appJwtSecret, err := client.GetVaultData(ctx, cfg.VaultSecret.App.Path, cfg.VaultSecret.App.JwtSecretVariable)
	//if err != nil {
	//	return err
	//}
	//if err = os.Setenv("APP_JWT_SECRET", appJwtSecret); err != nil {
	//	return err
	//}
	//
	//postgres_username, err := client.GetVaultData(ctx, cfg.VaultSecret.PostgresVault.Path, cfg.VaultSecret.PostgresVault.UsernameVariable)
	//if err != nil {
	//	return err
	//}
	//if err = os.Setenv("PG_USER", postgres_username); err != nil {
	//	return err
	//}
	//
	//postgres_password, err := client.GetVaultData(ctx, cfg.VaultSecret.PostgresVault.Path, cfg.VaultSecret.PostgresVault.PasswordVariable)
	//if err != nil {
	//	return err
	//}
	//if err = os.Setenv("PG_PASSWORD", postgres_password); err != nil {
	//	return err
	//}
	//
	//telebot_token, err := client.GetVaultData(ctx, cfg.VaultSecret.Telebot.Path, cfg.VaultSecret.Telebot.TokenVariable)
	//if err != nil {
	//	return err
	//}
	//if err = os.Setenv("TELEBOT_TOKEN", telebot_token); err != nil {
	//	return err
	//}

	fmt.Println(os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("APP_JWT_SECRET"), os.Getenv("TELEBOT_TOKEN"))

	return nil
}
