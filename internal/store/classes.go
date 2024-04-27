package store

import (
	"context"
	"fmt"
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"github.com/jackc/pgx/v5"
)

func (s *RgStore) CreateClass(class requests.CreateClass) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	var (
		query = `select * from classes.if_class_available($1, $2, $3)`
		count int
		id    = -1
		err   error
	)

	if err = s.s.QueryRow(ctx, query, class.Coach, class.ClassDate, class.ClassTime).Scan(&count); err != nil {
		return 0, err
	}

	if count > 0 {
		return id, fmt.Errorf("class already exists")
	}

	query = `select * from classes.create_class_if_not_exists($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	err = s.s.QueryRow(ctx, query,
		time.Now().UnixMilli(),
		class.Pupil,
		class.Coach,
		class.ClassDate,
		class.ClassTime,
		class.Duration,
		class.Price,
		class.ClassType,
		class.Capacity,
		class.IsOpen,
	).Scan(
		&id,
	)

	return id, err
}

func (s *RgStore) CancelClass(classId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from classes.cancel_class($1)`

	_, err := s.s.Exec(ctx, query, classId)

	return err
}

func (s *RgStore) DeleteClass(classId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from classes.delete_class($1)`

	_, err := s.s.Exec(ctx, query, classId)

	return err
}

func (s *RgStore) GetClassesForMonth(userType, today, lastDay string) ([]models.MicroClassInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()
	fmt.Println(today, lastDay)
	var (
		query   string
		class   models.MicroClassInfo
		classes []models.MicroClassInfo
	)

	switch userType {
	case "ADMIN":
		query = `select * from classes.get_classes_for_month_admin($1, $2)`
	}

	rows, err := s.s.Query(ctx, query, today, lastDay)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if _, err = pgx.ForEachRow(
		rows,
		[]any{
			&class.Key,
			&class.ClassDate,
			&class.ClassTime,
			&class.ClassDuration,
		},
		func() error {
			classes = append(classes, class)
			return nil
		}); err != nil {
		return nil, err
	}

	return classes, nil
}

func (s *RgStore) GetClassById(classId int) (*models.GetClassAdmin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	var class models.GetClassAdmin

	const query = `select * from classes.get_class_info($1)`

	err := s.s.QueryRow(ctx, query,
		classId,
	).Scan(
		&class.Key,
		&class.Pupils,
		&class.Coach,
		&class.ClassDate,
		&class.ClassTime,
		&class.ClassDuration,
		&class.ClassType,
		&class.Capacity,
		&class.Scheduled,
		&class.IsOpenToSignUp,
	)

	if err != nil {
		return nil, err
	}

	return &class, nil
}

func (s *RgStore) HaveAccessToClass(userID int) (bool, error) {
	_, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	// Получается здесь надо проверить что айдишник либо принадлежит тренеру и он
	//	имеет такое занятие, либо этот айдишник - ученицы и она тоже имеет такое занятие
	return false, nil
}
