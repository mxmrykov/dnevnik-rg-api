package repository

import (
	"context"
	postgres_requests "dnevnik-rg.ru/internal/postgres-requests"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{Pool: pool}
}

func (r *Repository) InitTablePupils() error {
	_, errInitTable := r.Pool.Exec(
		context.Background(),
		postgres_requests.InitTablePupils,
	)
	return errInitTable
}

func (r *Repository) InitTableCoaches() error {
	_, errInitTable := r.Pool.Exec(
		context.Background(),
		postgres_requests.InitTableCoaches,
	)
	return errInitTable
}

func (r *Repository) InitTablePasswords() error {
	_, errInitTable := r.Pool.Exec(
		context.Background(),
		postgres_requests.InitTablePasswords,
	)
	return errInitTable
}

func (r *Repository) InitTableClasses() error {
	_, errInitTable := r.Pool.Exec(
		context.Background(),
		postgres_requests.InitTableClasses,
	)
	return errInitTable
}
