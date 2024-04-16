package store

import (
	"context"

	"dnevnik-rg.ru/internal/models"
	"dnevnik-rg.ru/internal/models/response"
	"github.com/jackc/pgx/v5"
)

func (s *RgStore) NewAdmin(admin models.Admin) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from users.create_admin($1, $2, $3, $4, $5)`

	_, err := s.s.Exec(ctx, query, admin.Key, admin.Fio, admin.DateReg, admin.LogoUri, admin.Role)

	return err
}

func (s *RgStore) DeleteAdmin(key int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.delete_admin($1)`

	_, err := s.s.Exec(ctx, query, key)

	return err
}

func (s *RgStore) IsAdminExists(key int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select ex from users.if_admin_exists($1)`

	var exists bool

	err := s.s.QueryRow(ctx, query, key).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *RgStore) GetAdmin(key int) (*response.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.get_admin($1)`

	var admin response.Admin

	err := s.s.QueryRow(ctx, query, key).Scan(
		&admin.Key,
		&admin.Fio,
		&admin.DateReg,
		&admin.LogoUri,
		&admin.Role,
		&admin.Private.CheckSum,
		&admin.Private.LastUpdate,
		&admin.Private.Token,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (s *RgStore) GetAllAdmins() ([]models.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.list_admins()`
	var (
		admin  models.Admin
		admins []models.Admin
	)

	rows, _ := s.s.Query(ctx, query)

	_, err := pgx.ForEachRow(
		rows,
		[]any{
			nil,
			&admin.Key,
			&admin.Fio,
			&admin.DateReg,
			&admin.LogoUri,
			&admin.Role,
		},
		func() error {
			admins = append(admins, admin)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (s *RgStore) GetAllAdminsExcept(key int) ([]models.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from users.list_admins_except($1)`
	var (
		admin  models.Admin
		admins []models.Admin
	)

	rows, _ := s.s.Query(ctx, query, key)

	_, err := pgx.ForEachRow(
		rows,
		[]any{
			&admin.Key,
			&admin.Fio,
			&admin.DateReg,
			&admin.LogoUri,
			&admin.Role,
		},
		func() error {
			admins = append(admins, admin)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (s *RgStore) GetAdminClassesForToday(date string) ([]models.ShortClassInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from classes.get_classes_for_today_admin($1)`
	var (
		classInfo   models.ShortClassInfo
		classesInfo []models.ShortClassInfo
	)

	rows, _ := s.s.Query(ctx, query, date)

	_, err := pgx.ForEachRow(
		rows,
		[]any{
			nil,
			&classInfo.Key,
			&classInfo.Pupils,
			&classInfo.Coach,
			&classInfo.ClassTime,
			&classInfo.ClassDuration,
			&classInfo.ClassType,
			&classInfo.Scheduled,
			&classInfo.PupilCount,
		},
		func() error {
			classesInfo = append(classesInfo, classInfo)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return classesInfo, nil
}
