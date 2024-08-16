package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type (
	Config struct {
		App         `yaml:"app"`
		Http        `yaml:"http"`
		Postgres    `yaml:"postgres"`
		Telebot     `yaml:"telebot"`
		Vault       `yaml:"vault_config"`
		VaultSecret `yaml:"vault"`
	}

	App struct {
		Name      string `yaml:"app_name" env:"APP_NAME"`
		Version   string `yaml:"app_version" env:"APP_VERSION"`
		Deploy    string `yaml:"deploy" env:"DEPLOY"`
		JwtSecret string `required:"true" yaml:"jwt_secret" env:"APP_JWT_SECRET"`
	}

	Http struct {
		Host string `required:"true" yaml:"http_host" env:"HTTP_HOST"`
		Port string `required:"true" yaml:"http_port" env:"HTTP_PORT"`
	}

	Postgres struct {
		Host     string `required:"true" yaml:"postgres_host" env:"PG_HOST"`
		Port     string `required:"true" yaml:"postgres_port" env:"PG_PORT"`
		User     string `required:"true" yaml:"postgres_username" env:"PG_USER"`
		Password string `required:"true" yaml:"postgres_password" env:"PG_PASSWORD"`
		DBName   string `required:"true" yaml:"postgres_dbname" env:"PG_NAME"`
	}

	Vault struct {
		Host string `required:"true" yaml:"vault_host" env:"VAULT_HOST"`
		Port string `required:"true" yaml:"vault_port" env:"VAULT_PORT"`
	}

	Telebot struct {
		Token string `required:"true" yaml:"token" env:"TELEBOT_TOKEN"`
	}

	VaultSecret struct {
		PostgresVault struct {
			Path             string `yaml:"path"`
			UsernameVariable string `yaml:"username_variable"`
			PasswordVariable string `yaml:"password_variable"`
		} `yaml:"postgres"`
		App struct {
			Path              string `yaml:"path"`
			JwtSecretVariable string `yaml:"jwt_secret_variable"`
		} `yaml:"app"`
		Telebot struct {
			Path          string `yaml:"path"`
			TokenVariable string `yaml:"token_variable"`
		} `yaml:"telebot"`
	}

	// VaultCfg especially for vault extractor
	VaultCfg struct {
		Vault       `yaml:"vault_config"`
		VaultSecret `yaml:"vault"`
	}
)

func NewConfig() (*Config, error) {
	config := new(Config)
	if errReadConfig := cleanenv.ReadConfig(getConfigPath(), config); errReadConfig != nil {
		return nil, fmt.Errorf("error read config: %v", errReadConfig)
	}
	if err := os.Setenv("JWT_SECRET", config.JwtSecret); err != nil {
		return nil, err
	}
	if err := os.Setenv("APP_NAME", config.Name); err != nil {
		return nil, err
	}
	if err := os.Setenv("APP_VERSION", config.Version); err != nil {
		return nil, err
	}
	if err := os.Setenv("DEPLOY", config.Deploy); err != nil {
		return nil, err
	}
	return config, nil
}

func NewVaultConfig() (*VaultCfg, error) {
	vaultConfig := new(VaultCfg)

	if err := cleanenv.ReadConfig(getConfigPath(), vaultConfig); err != nil {
		return nil, err
	}

	return vaultConfig, nil
}

func getConfigPath() string {
	path := "config/local/config.yml"
	log.Info().Msg(os.Getenv("BUILD_ENV") + "," + os.Getenv("VAULT_DEV_ROOT_TOKEN_ID"))
	if os.Getenv("BUILD_ENV") == "stage" {
		path = "/bin/stage/config.yml"
	} else if os.Getenv("BUILD_ENV") == "prod" {
		path = "/bin/prod/config.yml"
	}

	return path
}
