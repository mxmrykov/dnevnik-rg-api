package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/internal/models/response"
	pgRequests "dnevnik-rg.ru/internal/postgres-requests"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	Shard1 *pgxpool.Pool
	Shard2 *pgxpool.Pool
	Shard3 *pgxpool.Pool
}

func NewRepository(shards []*pgxpool.Pool) *Repository {
	return &Repository{
		Shard1: shards[0],
		Shard2: shards[1],
		Shard3: shards[2],
	}
}

func (r *Repository) InitTablePupils() error {
	_, errInitTable := r.Shard1.Exec(
		context.Background(),
		pgRequests.InitTablePupils,
	)
	return errInitTable
}

func (r *Repository) InitTableCoaches() error {
	_, errInitTable := r.Shard1.Exec(
		context.Background(),
		pgRequests.InitTableCoaches,
	)
	return errInitTable
}

func (r *Repository) InitTablePasswords() error {
	_, errInitTable := r.Shard2.Exec(
		context.Background(),
		pgRequests.InitTablePasswords,
	)
	return errInitTable
}

func (r *Repository) InitTableClasses() error {
	_, errInitTable := r.Shard3.Exec(
		context.Background(),
		pgRequests.InitTableClasses,
	)
	return errInitTable
}

func (r *Repository) InitTableAdmins() error {
	_, errInitTable := r.Shard1.Exec(
		context.Background(),
		pgRequests.InitTableAdmins,
	)
	return errInitTable
}

func (r *Repository) GetAllPupils() ([]models.Pupil, error) {
	var pupils []models.Pupil
	rows, err := r.Shard1.Query(
		context.Background(),
		pgRequests.GetAllPupils,
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
	rows, err := r.Shard1.Query(
		context.Background(),
		pgRequests.GetAllCoaches,
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
	rows, err := r.Shard1.Query(
		context.Background(),
		pgRequests.GetAllAdmins,
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
	rows, err := r.Shard1.Query(
		context.Background(),
		pgRequests.GetAllAdminsExcept,
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
	_, errNewAdmin := r.Shard1.Exec(
		context.Background(),
		pgRequests.NewAdmin,
		admin.Key, admin.Fio, admin.DateReg, admin.LogoUri, "ADMIN",
	)
	return errNewAdmin
}

func (r *Repository) DeleteAdmin(key int) error {
	_, errNewAdmin := r.Shard1.Exec(
		context.Background(),
		pgRequests.DeleteAdmin,
		key,
	)
	return errNewAdmin
}

func (r *Repository) NewPassword(password models.Password) error {
	_, errNewAdmin := r.Shard2.Exec(
		context.Background(),
		pgRequests.NewPassword,
		password.Key, password.CheckSum, password.Token, password.LastUpdate,
	)
	return errNewAdmin
}

func (r *Repository) DeletePassword(key int) error {
	_, errNewAdmin := r.Shard2.Exec(
		context.Background(),
		pgRequests.DeletePassword,
		key,
	)
	return errNewAdmin
}

func (r *Repository) GetAdmin(key int) (response.Admin, error) {
	var Admin response.Admin
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.GetAdmin,
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
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.IsAdminExists,
		key,
	).Scan(
		&count,
	)
	return count == 1, err
}

func (r *Repository) CreateCoach(coach models.Coach) error {
	_, errNewCoach := r.Shard1.Exec(
		context.Background(),
		pgRequests.CreateCoach,
		coach.Key, coach.Fio, coach.DateReg,
		coach.HomeCity, coach.TrainingCity,
		coach.Birthday, coach.About, coach.LogoUri,
		"COACH",
	)
	return errNewCoach
}

func (r *Repository) GetCoach(key int) (response.Coach, error) {
	var coach response.Coach
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.GetCoach,
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
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.GetCoachFull,
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
	_, err := r.Shard1.Exec(
		context.Background(),
		sql,
	)
	return err
}

func (r *Repository) DeleteCoach(key int) error {
	_, err := r.Shard1.Exec(
		context.Background(),
		pgRequests.DeleteCoach,
		key,
	)
	return err
}

func (r *Repository) GetCoachPupils(coachId int) ([]response.PupilList, error) {
	var pupils []response.PupilList
	rows, err := r.Shard1.Query(
		context.Background(),
		pgRequests.GetCoachPupils,
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
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.IsCoachExists,
		key,
	).Scan(
		&count,
	)
	return count == 1, err
}

func (r *Repository) CreatePupil(Pupil models.Pupil) error {
	_, errNewCoach := r.Shard1.Exec(
		context.Background(),
		pgRequests.CreatePupil,
		Pupil.Key, Pupil.Fio, Pupil.DateReg,
		Pupil.Coach, Pupil.HomeCity, Pupil.TrainingCity,
		Pupil.Birthday, Pupil.About, Pupil.CoachReview,
		Pupil.LogoUri, "PUPIL",
	)
	return errNewCoach
}

func (r *Repository) GetPupilFull(key int) (response.PupilFull, error) {
	var pupil response.PupilFull
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.GetPupilFull,
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
	err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.GetPupil,
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
	_, err := r.Shard1.Exec(
		context.Background(),
		sql,
	)
	return err
}

func (r *Repository) DeletePupil(key int) error {
	_, err := r.Shard1.Exec(
		context.Background(),
		pgRequests.DeletePupil,
		key,
	)
	return err
}

func (r *Repository) Authorize(key int, checksum string) (response.Auth, error) {
	var (
		password string
		auth     response.Auth
	)

	if err := r.Shard2.QueryRow(
		context.Background(),
		pgRequests.GetPasswordCheck,
		key,
	).Scan(
		&password,
		&auth.Token,
	); err != nil {
		return response.Auth{}, err
	}

	if password != checksum {
		return response.Auth{}, fmt.Errorf("wrong password")
	}

	if err := r.Shard1.QueryRow(
		context.Background(),
		pgRequests.Auth,
		key,
	).Scan(
		&auth.Role,
	); err != nil {
		return response.Auth{}, err
	}

	auth.Key = key

	return auth, nil
}

func (r *Repository) GetBirthdaysList(key int) ([]requests.BirthDayList, error) {
	var bdays []requests.BirthDayList
	rows, err := r.Shard1.Query(
		context.Background(),
		pgRequests.GetCoachNearestBirthdays,
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

	rows, err := r.Shard3.Query(
		context.Background(),
		pgRequests.GetCoachSchedule,
		key,
		date,
	)

	log.Println(key, date)
	for rows.Next() {
		var class models.ClassMainInfo
		if errScan := rows.Scan(
			&class.Key,
			&class.Pupil,
			&class.Coach,
			&class.ClassDate,
			&class.ClassTime,
			&class.ClassDuration,
		); errScan != nil {
			log.Println("error while scanning main class info: ", errScan)
		}
		log.Println(class)
		classes = append(classes, class)
	}
	if err != nil {
		return nil, err
	}
	log.Println(classes)
	return classes, nil
}

func (r *Repository) AutoUpdateToken(token string, key int) error {
	_, err := r.Shard2.Exec(
		context.Background(),
		pgRequests.UpdateOldToken,
		key, token, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *Repository) CreateClass(class requests.CreateClass) (id int, err error) {
	var count int8
	err = r.Shard3.QueryRow(
		context.Background(),
		pgRequests.IfClassAvail, class.Coach,
		class.ClassDate,
		class.ClassTime,
	).Scan(
		&count,
	)
	if count > 0 {
		id = -1

		return
	}
	err = r.Shard3.QueryRow(
		context.Background(),
		pgRequests.CreateClassIfNotExists,
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

	return
}
