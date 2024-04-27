package postgres

import (
	"context"
	"fmt"

	"dnevnik-rg.ru/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(postgresConfig *config.Postgres) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresConfig.Shard1.User,
		postgresConfig.Shard1.Password,
		postgresConfig.Shard1.Host,
		postgresConfig.Shard1.Port,
		postgresConfig.Shard1.DBName,
	))
}
