package postgres

import (
	"context"
	"fmt"

	"dnevnik-rg.ru/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgres(postgresConfig *config.Postgres) ([]*pgxpool.Pool, error) {
	var shards []*pgxpool.Pool

	shard, err := pgxpool.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresConfig.Shard1.User,
		postgresConfig.Shard1.Password,
		postgresConfig.Shard1.Host,
		postgresConfig.Shard1.Port,
		postgresConfig.Shard1.DBName,
	))
	if err != nil {
		return nil, fmt.Errorf("error at connecting to shard 1: %v", err)
	}
	shards = append(shards, shard)
	shard, err = pgxpool.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresConfig.Shard2.User,
		postgresConfig.Shard2.Password,
		postgresConfig.Shard2.Host,
		postgresConfig.Shard2.Port,
		postgresConfig.Shard2.DBName,
	))
	if err != nil {
		return nil, fmt.Errorf("error at connecting to shard 2: %v", err)
	}
	shards = append(shards, shard)
	shard, err = pgxpool.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresConfig.Shard3.User,
		postgresConfig.Shard3.Password,
		postgresConfig.Shard3.Host,
		postgresConfig.Shard3.Port,
		postgresConfig.Shard3.DBName,
	))
	if err != nil {
		return nil, fmt.Errorf("error at connecting to shard 3: %v", err)
	}
	shards = append(shards, shard)

	return shards, nil
}
