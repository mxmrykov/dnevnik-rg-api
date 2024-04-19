package store

import (
	"context"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/internal/models/response"
	"github.com/jackc/pgx/v5"
)

func (s *RgStore) GetAllCoaches() ([]models.Coach, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_all_coaches()`
	var (
		coach   models.Coach
		coaches []models.Coach
	)

	rows, _ := s.s.Query(ctx, query)

	_, err := pgx.ForEachRow(
		rows,
		[]any{
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
		},
		func() error {
			coaches = append(coaches, coach)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return coaches, nil
}

func (s *RgStore) CreateCoach(coach models.Coach) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from users.create_coach($1, $2, $3, $4, $5, $6, $7, $8, 'COACH')`

	_, err := s.s.Exec(ctx, query, coach.Key, coach.Fio, coach.DateReg,
		coach.HomeCity, coach.TrainingCity,
		coach.Birthday, coach.About, coach.LogoUri)

	return err
}

func (s *RgStore) GetCoach(key int) (*response.Coach, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_coach($1)`

	var coach response.Coach

	err := s.s.QueryRow(ctx, query, key).Scan(
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

	if err != nil {
		return nil, err
	}

	return &coach, nil
}

func (s *RgStore) GetCoachFull(key int) (*response.CoachFull, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_coach_full($1)`

	var coach response.CoachFull

	err := s.s.QueryRow(ctx, query, key).Scan(
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

	if err != nil {
		return nil, err
	}

	return &coach, nil
}

func (s *RgStore) UpdateCoach(sql string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	_, err := s.s.Exec(ctx, sql)

	return err
}

func (s *RgStore) DeleteCoach(key int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from users.delete_coach($1)`

	_, err := s.s.Exec(ctx, query, key)

	return err
}

func (s *RgStore) GetCoachPupils(coachId int) ([]response.PupilList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_coach_pupils($1)`
	var (
		pupil  response.PupilList
		pupils []response.PupilList
	)

	rows, err := s.s.Query(ctx, query, coachId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	_, err = pgx.ForEachRow(
		rows,
		[]any{
			&pupil.Key,
			&pupil.Fio,
			&pupil.LogoUri,
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

func (s *RgStore) IsCoachExists(key int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select count from users.if_coach_exists($1)`

	var count int

	err := s.s.QueryRow(ctx, query, key).Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (s *RgStore) GetBirthdaysList(key int) ([]requests.BirthDayList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_nearest_pupils_bd($1)`
	var (
		bday  requests.BirthDayList
		bdays []requests.BirthDayList
	)

	rows, err := s.s.Query(ctx, query, key)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	_, err = pgx.ForEachRow(
		rows,
		[]any{
			&bday.Key,
			&bday.Fio,
			&bday.Birthday,
		},
		func() error {
			bdays = append(bdays, bday)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return bdays, nil
}

func (s *RgStore) GetCoachSchedule(key int, date string) ([]models.ClassMainInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from classes.get_coach_schedule($1, $2)`
	var (
		class   models.ClassMainInfo
		classes []models.ClassMainInfo
	)

	rows, err := s.s.Query(ctx, query, key, date)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	_, err = pgx.ForEachRow(
		rows,
		[]any{
			&class.Key,
			&class.Pupils,
			&class.Coach,
			&class.ClassDate,
			&class.ClassTime,
			&class.ClassDuration,
		},
		func() error {
			classes = append(classes, class)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return classes, nil
}

func (s *RgStore) GetCoachNameById(id int) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select fio from users.coaches where key = $1`

	var fio string

	if err := s.s.QueryRow(ctx, query, id).Scan(&fio); err != nil {
		return nil, err
	}

	return &fio, nil
}
