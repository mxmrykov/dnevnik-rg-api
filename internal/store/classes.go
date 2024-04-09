package store

import (
	"context"
	"fmt"
	"time"

	requests "dnevnik-rg.ru/internal/models/request"
)

func (s *RgStore) CreateClass(class requests.CreateClass) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	var (
		query = `select * from classes.if_class_available($1, $2, $3)`
		count int
		id    = -1
	)

	err := s.s.QueryRow(ctx, query, class.Coach, class.ClassDate, class.ClassTime).Scan(&count)

	if err != nil {
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

	return id, nil
}

func (s *RgStore) CancelClass(classId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from classes.cancel_class($1)`

	_, err := s.s.Exec(ctx, query, classId)

	return err
}
