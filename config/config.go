package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"APP"`
		Http     `yaml:"HTTP"`
		Postgres `yaml:"POSTGRES"`
		TgBots   `yaml:"TG_BOTS"`
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
		Shard1 struct {
			Host     string `required:"true" yaml:"PG_HOST" env:"PG_HOST"`
			Port     string `required:"true" yaml:"PG_PORT" env:"PG_PORT"`
			User     string `required:"true" yaml:"PG_USER" env:"PG_USER"`
			Password string `required:"true" yaml:"PG_PASSWORD" env:"PG_PASSWORD"`
			DBName   string `required:"true" yaml:"PG_NAME" env:"PG_NAME"`
			PgDriver string `required:"true" yaml:"PG_DRIVER" env:"PG_DRIVER"`
		} `required:"true" yaml:"shard1" env:"shard1"`
		Shard2 struct {
			Host     string `required:"true" yaml:"PG_HOST" env:"PG_HOST"`
			Port     string `required:"true" yaml:"PG_PORT" env:"PG_PORT"`
			User     string `required:"true" yaml:"PG_USER" env:"PG_USER"`
			Password string `required:"true" yaml:"PG_PASSWORD" env:"PG_PASSWORD"`
			DBName   string `required:"true" yaml:"PG_NAME" env:"PG_NAME"`
			PgDriver string `required:"true" yaml:"PG_DRIVER" env:"PG_DRIVER"`
		} `required:"true" yaml:"shard2" env:"shard2"`
		Shard3 struct {
			Host     string `required:"true" yaml:"PG_HOST" env:"PG_HOST"`
			Port     string `required:"true" yaml:"PG_PORT" env:"PG_PORT"`
			User     string `required:"true" yaml:"PG_USER" env:"PG_USER"`
			Password string `required:"true" yaml:"PG_PASSWORD" env:"PG_PASSWORD"`
			DBName   string `required:"true" yaml:"PG_NAME" env:"PG_NAME"`
			PgDriver string `required:"true" yaml:"PG_DRIVER" env:"PG_DRIVER"`
		} `required:"true" yaml:"shard3" env:"shard3"`
	}
	TgBots struct {
		TgTechBot struct {
			Token string `required:"true" yaml:"TG_TECH_TOKEN" env:"TG_TECH_TOKEN"`
		} `required:"true" yaml:"TECH_BOT" env:"TECH_BOT"`
	}
)

func NewConfig() (*Config, error) {
	config := &Config{}
	// stageConfigPath := "/bin/config.yml"
	localConfigPath := "config/local/config.yml"
	if errReadConfig := cleanenv.ReadConfig(localConfigPath, config); errReadConfig != nil {
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
