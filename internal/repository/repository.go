package repository

import (
	"context"
	"log"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
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

func (r *Repository) GetAllPupils() ([]models.Pupil, error) {
	var pupils []models.Pupil
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetAllPupils,
	)
	for rows.Next() {
		var pupil models.Pupil
		if errScan := rows.Scan(
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
		); errScan != nil {
			log.Println("error while scanning recovery pupils from DB: ", errScan)
		}
		pupils = append(pupils, pupil)
	}
	if err != nil {
		return nil, err
	}
	return pupils, nil
}

func (r *Repository) GetAllCoaches() ([]models.Coach, error) {
	var coaches []models.Coach
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetAllCoaches,
	)
	for rows.Next() {
		var coach models.Coach
		if errScan := rows.Scan(
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
		); errScan != nil {
			log.Println("error while scanning recovery coaches from DB: ", errScan)
		}
		coaches = append(coaches, coach)
	}
	if err != nil {
		return nil, err
	}
	return coaches, nil
}

func (r *Repository) GetAllAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetAllAdmins,
	)
	for rows.Next() {
		var admin models.Admin
		if errScan := rows.Scan(
			nil,
			&admin.Key,
			&admin.Fio,
			&admin.DateReg,
			&admin.LogoUri,
			&admin.Role,
		); errScan != nil {
			log.Println("error while scanning recovery admins from DB: ", errScan)
		}
		admins = append(admins, admin)
	}
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *Repository) GetAllAdminsExcept(key int) ([]response.AdminList, error) {
	var admins []response.AdminList
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetAllAdminsExcept,
		key,
	)
	for rows.Next() {
		var admin response.AdminList
		if errScan := rows.Scan(
			&admin.Key,
			&admin.Fio,
			&admin.LogoUri,
		); errScan != nil {
			log.Println("error while scanning list admins from db: ", errScan)
		}
		admins = append(admins, admin)
	}
	if err != nil {
		return nil, err
	}
	return admins, nil
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
		&Admin.Role,
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

func (r *Repository) UpdateCoach(sql string) error {
	_, err := r.Pool.Exec(
		context.Background(),
		sql,
	)
	return err
}

func (r *Repository) DeleteCoach(key int) error {
	_, err := r.Pool.Exec(
		context.Background(),
		postgres_requests.DeleteCoach,
		key,
	)
	return err
}

func (r *Repository) GetCoachPupils(coachId int) ([]response.PupilList, error) {
	var pupils []response.PupilList
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetCoachPupils,
		coachId,
	)
	for rows.Next() {
		var pupil response.PupilList
		if errScan := rows.Scan(
			&pupil.Key,
			&pupil.Fio,
			&pupil.LogoUri,
		); errScan != nil {
			log.Println("error while scanning recovery pupils from DB: ", errScan)
		}
		pupils = append(pupils, pupil)
	}
	if err != nil {
		return nil, err
	}
	return pupils, nil
}

func (r *Repository) IsCoachExists(key int) (bool, error) {
	var count int
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.IsCoachExists,
		key,
	).Scan(
		&count,
	)
	return count == 1, err
}

func (r *Repository) CreatePupil(Pupil models.Pupil) error {
	_, errNewCoach := r.Pool.Exec(
		context.Background(),
		postgres_requests.CreatePupil,
		Pupil.Key, Pupil.Fio, Pupil.DateReg,
		Pupil.Coach, Pupil.HomeCity, Pupil.TrainingCity,
		Pupil.Birthday, Pupil.About, Pupil.CoachReview,
		Pupil.LogoUri, "PUPIL",
	)
	return errNewCoach
}

func (r *Repository) GetPupilFull(key int) (response.PupilFull, error) {
	var pupil response.PupilFull
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.GetPupilFull,
		key,
	).Scan(
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
		&pupil.Private.CheckSum,
		&pupil.Private.Token,
		&pupil.Private.LastUpdate,
	)
	return pupil, err
}

func (r *Repository) GetPupil(key int) (response.Pupil, error) {
	var pupil response.Pupil
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.GetPupil,
		key,
	).Scan(
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
	return pupil, err
}

func (r *Repository) UpdatePupil(sql string) error {
	_, err := r.Pool.Exec(
		context.Background(),
		sql,
	)
	return err
}

func (r *Repository) DeletePupil(key int) error {
	_, err := r.Pool.Exec(
		context.Background(),
		postgres_requests.DeletePupil,
		key,
	)
	return err
}

func (r *Repository) Authorize(key int, checksum string) (response.Auth, error) {
	var auth response.Auth
	err := r.Pool.QueryRow(
		context.Background(),
		postgres_requests.Auth,
		key, checksum,
	).Scan(
		&auth.Key,
		&auth.Token,
		&auth.Role,
	)
	return auth, err
}

func (r *Repository) GetBirthdaysList(key int) ([]requests.BirthDayList, error) {
	var bdays []requests.BirthDayList
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetCoachNearestBirthdays,
		key,
	)
	for rows.Next() {
		var bday requests.BirthDayList
		if errScan := rows.Scan(
			&bday.Key,
			&bday.Fio,
			&bday.Birthday,
		); errScan != nil {
			log.Println("error while scanning recovery pupils from DB: ", errScan)
		}
		bdays = append(bdays, bday)
	}
	if err != nil {
		return nil, err
	}
	return bdays, nil
}

func (r *Repository) GetCoachSchedule(key int, date string) ([]models.ClassMainInfo, error) {
	var classes []models.ClassMainInfo
	rows, err := r.Pool.Query(
		context.Background(),
		postgres_requests.GetCoachSchedule,
		key,
		date,
	)
	for rows.Next() {
		var class models.ClassMainInfo
		if errScan := rows.Scan(
			nil,
			&class.Key,
			&class.Pupil,
			&class.Coach,
			&class.ClassDate,
			&class.ClassTime,
			&class.ClassDuration,
		); errScan != nil {
			log.Println("error while scanning main class info: ", errScan)
		}
		classes = append(classes, class)
	}
	if err != nil {
		return nil, err
	}
	return classes, nil
}
