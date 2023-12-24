package postgres

import (
	"context"
	"dnevnik-rg.ru/config"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgres(postgresConfig *config.Postgres) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.DBName,
	))
}
