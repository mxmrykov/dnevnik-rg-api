package store

import (
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/internal/models/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	GetAllPupils() ([]models.Pupil, error)
	CreatePupil(Pupil models.Pupil) error
	GetPupilFull(key int) (*response.PupilFull, error)
	GetPupil(key int) (*response.Pupil, error)
	UpdatePupil(sql string) error
	DeletePupil(key int) error
	GetPupilsNameByIds(ids []int) ([]string, error)

	NewAdmin(admin models.Admin) error
	DeleteAdmin(key int) error
	GetAllAdmins() ([]models.Admin, error)
	GetAllAdminsExcept(key int) ([]models.Admin, error)
	IsAdminExists(key int) (bool, error)
	GetAdmin(key int) (*response.Admin, error)
	GetAdminClassesForToday(date string) ([]models.ShortClassInfo, error)

	GetAllCoaches() ([]models.Coach, error)
	CreateCoach(coach models.Coach) error
	GetCoach(key int) (*response.Coach, error)
	GetCoachFull(key int) (*response.CoachFull, error)
	UpdateCoach(sql string) error
	DeleteCoach(key int) error
	GetCoachPupils(coachId int) ([]response.PupilList, error)
	IsCoachExists(key int) (bool, error)
	GetBirthdaysList(key int) ([]requests.BirthDayList, error)
	GetCoachSchedule(key int, date string) ([]models.ClassMainInfo, error)
	GetCoachNameById(id int) (*string, error)

	NewPassword(password models.Password) error
	DeletePassword(key int) error
	AutoUpdateToken(token string, key int) error

	Authorize(key int, checksum string, ip string) (*response.Auth, error)

	CreateClass(class requests.CreateClass) (id int, err error)
	CancelClass(classId int) error
	DeleteClass(classId int) error
}

type RgStore struct {
	s                *pgxpool.Pool
	operationTimeout time.Duration
}

func NewStore(store *pgxpool.Pool, timeOut time.Duration) *RgStore {
	return &RgStore{
		s:                store,
		operationTimeout: timeOut,
	}
}
