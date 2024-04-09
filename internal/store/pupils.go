package store

import (
	"context"

	"dnevnik-rg.ru/internal/models"
	"dnevnik-rg.ru/internal/models/response"
	"github.com/jackc/pgx/v5"
)

func (s *RgStore) GetAllPupils() ([]models.Pupil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_all_pupils()`
	var (
		pupil  models.Pupil
		pupils []models.Pupil
	)

	rows, _ := s.s.Query(ctx, query)

	_, err := pgx.ForEachRow(
		rows,
		[]any{
			&pupil.Key,
			&pupil.Fio,
			&pupil.DateReg,
			&pupil.Coach,
			&pupil.HomeCity,
			&pupil.TrainingCity,
			&pupil.Birthday,
			&pupil.About,
			&pupil.CoachReview,
			&pupil.LogoUri,
			&pupil.Role,
		},
		func() error {
			pupils = append(pupils, pupil)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return pupils, nil
}

func (s *RgStore) CreatePupil(Pupil models.Pupil) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from users.create_pupil($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'PUPIL')`

	_, err := s.s.Exec(ctx, query, Pupil.Key, Pupil.Fio, Pupil.DateReg,
		Pupil.Coach, Pupil.HomeCity, Pupil.TrainingCity,
		Pupil.Birthday, Pupil.About, Pupil.CoachReview,
		Pupil.LogoUri)

	return err
}

func (s *RgStore) GetPupilFull(key int) (*response.PupilFull, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_pupil_full($1)`

	var pupil response.PupilFull

	err := s.s.QueryRow(ctx, query, key).Scan(
		&pupil.Key,
		&pupil.Fio,
		&pupil.DateReg,
		&pupil.HomeCity,
		&pupil.TrainingCity,
		&pupil.Birthday,
		&pupil.About,
		&pupil.LogoUri,
		&pupil.Role,
		&pupil.Private.CheckSum,
		&pupil.Private.Token,
		&pupil.Private.LastUpdate,
	)

	if err != nil {
		return nil, err
	}

	return &pupil, nil
}

func (s *RgStore) GetPupil(key int) (*response.Pupil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_pupil($1)`

	var pupil response.Pupil

	err := s.s.QueryRow(ctx, query, key).Scan(
		nil,
		&pupil.Key,
		&pupil.Fio,
		&pupil.DateReg,
		&pupil.Coach,
		&pupil.HomeCity,
		&pupil.TrainingCity,
		&pupil.Birthday,
		&pupil.About,
		&pupil.CoachReview,
		&pupil.LogoUri,
		&pupil.Role,
	)

	if err != nil {
		return nil, err
	}

	return &pupil, nil
}

func (s *RgStore) UpdatePupil(sql string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	_, err := s.s.Exec(ctx, sql)

	return err
}

func (s *RgStore) DeletePupil(key int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from users.delete_pupil($1)`

	_, err := s.s.Exec(ctx, query, key)

	return err
}

func (s *RgStore) GetPupilsNameByIds(ids []int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_pupils_names($1)`
	var (
		name  string
		names []string
	)

	rows, _ := s.s.Query(ctx, query, ids)

	_, err := pgx.ForEachRow(
		rows,
		[]any{
			&name,
		},
		func() error {
			names = append(names, name)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return names, nil
}
