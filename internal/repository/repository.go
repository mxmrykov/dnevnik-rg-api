package repository

import (
	"context"
	"dnevnik-rg.ru/internal/models"
	"dnevnik-rg.ru/internal/models/response"
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

func (r *Repository) InitTableAdmins() error {
	_, errInitTable := r.Pool.Exec(
		context.Background(),
		postgres_requests.InitTableAdmins,
	)
	return errInitTable
}

func (r *Repository) NewAdmin(admin models.Admin) error {
	_, errNewAdmin := r.Pool.Exec(
		context.Background(),
		postgres_requests.NewAdmin,
		admin.Key, admin.Fio, admin.DateReg, admin.LogoUri, "ADMIN",
	)
	return errNewAdmin
}

func (r *Repository) DeleteAdmin(key int) error {
	_, errNewAdmin := r.Pool.Exec(
		context.Background(),
		postgres_requests.DeleteAdmin,
		key,
	)
	return errNewAdmin
}

func (r *Repository) NewPassword(password models.Password) error {
	_, errNewAdmin := r.Pool.Exec(
		context.Background(),
		postgres_requests.NewPassword,
		password.Key, password.CheckSum, password.Token, password.LastUpdate,
	)
	return errNewAdmin
}

func (r *Repository) DeletePassword(key int) error {
	_, errNewAdmin := r.Pool.Exec(
		context.Background(),
		postgres_requests.DeletePassword,
		key,
	)
	return errNewAdmin
}

func (r *Repository) GetAdmin(key int) (response.Admin, error) {
	var Admin response.Admin
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.GetAdmin,
		key,
	).Scan(
		&Admin.Key,
		&Admin.Fio,
		&Admin.DateReg,
		&Admin.LogoUri,
		&Admin.Private.CheckSum,
		&Admin.Private.LastUpdate,
		&Admin.Private.Token,
	)
	return Admin, err
}

func (r *Repository) IsAdminExists(key int) (bool, error) {
	var count int
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.IsAdminExists,
		key,
	).Scan(
		&count,
	)
	return count == 1, err
}

func (r *Repository) CreateCoach(coach models.Coach) error {
	_, errNewCoach := r.Pool.Exec(
		context.Background(),
		postgres_requests.CreateCoach,
		coach.Key, coach.Fio, coach.DateReg,
		coach.HomeCity, coach.TrainingCity,
		coach.Birthday, coach.About, coach.LogoUri,
		"COACH",
	)
	return errNewCoach
}

func (r *Repository) GetCoach(key int) (response.Coach, error) {
	var coach response.Coach
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.GetCoach,
		key,
	).Scan(
		nil,
		&coach.Key,
		&coach.Fio,
		&coach.DateReg,
		&coach.HomeCity,
		&coach.TrainingCity,
		&coach.Birthday,
		&coach.About,
		&coach.LogoUri,
		&coach.Role,
	)
	return coach, err
}

func (r *Repository) GetCoachFull(key int) (response.CoachFull, error) {
	var coach response.CoachFull
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.GetCoachFull,
		key,
	).Scan(
		&coach.Key,
		&coach.Fio,
		&coach.DateReg,
		&coach.HomeCity,
		&coach.TrainingCity,
		&coach.Birthday,
		&coach.About,
		&coach.LogoUri,
		&coach.Role,
		&coach.Private.CheckSum,
		&coach.Private.Token,
		&coach.Private.LastUpdate,
	)
	return coach, err
}

func (r *Repository) DeleteCoach(key int) error {
	_, err := r.Pool.Exec(
		context.Background(),
		postgres_requests.DeleteCoach,
		key,
	)
	return err
}
