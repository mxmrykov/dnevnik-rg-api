package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	Config struct {
		App      `yaml:"APP"`
		Http     `yaml:"HTTP"`
		Postgres `yaml:"POSTGRES"`
	}

	App struct {
		Name      string `yaml:"APP_NAME" env:"APP_NAME"`
		Version   string `yaml:"APP_VERSION" env:"APP_VERSION"`
		Deploy    string `yaml:"DEPLOY" env:"DEPLOY"`
		JwtSecret string `required:"true" yaml:"JWT_SECRET" env:"JWT_SECRET"`
	}

	Http struct {
		Host string `required:"true" yaml:"HTTP_HOST" env:"HTTP_HOST"`
		Port string `required:"true" yaml:"HTTP_PORT" env:"HTTP_PORT"`
	}

	Postgres struct {
		Host     string `required:"true" yaml:"PG_HOST" env:"PG_HOST"`
		Port     string `required:"true" yaml:"PG_PORT" env:"PG_PORT"`
		User     string `required:"true" yaml:"PG_USER" env:"PG_USER"`
		Password string `required:"true" yaml:"PG_PASSWORD" env:"PG_PASSWORD"`
		DBName   string `required:"true" yaml:"PG_NAME" env:"PG_NAME"`
		PgDriver string `required:"true" yaml:"PG_DRIVER" env:"PG_DRIVER"`
	}
)

func NewConfig() (*Config, error) {
	config := &Config{}
	stageConfigPath := "./config/stage/config.yml"
	if errReadConfig := cleanenv.ReadConfig(stageConfigPath, config); errReadConfig != nil {
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
