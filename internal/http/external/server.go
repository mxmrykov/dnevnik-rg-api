package external

import (
	"net/http"

	"dnevnik-rg.ru/internal/cache"
	"dnevnik-rg.ru/internal/models"
	"dnevnik-rg.ru/internal/repository"
)

type Server interface {
	CreateAdmin(write http.ResponseWriter, request *http.Request)
	GetAdmin(write http.ResponseWriter, request *http.Request)
	GetAllAdminsExcept(write http.ResponseWriter, request *http.Request)

	CreateCoach(write http.ResponseWriter, request *http.Request)
	GetCoach(write http.ResponseWriter, request *http.Request)
	GetCoachFull(write http.ResponseWriter, request *http.Request)
	UpdateCoach(write http.ResponseWriter, request *http.Request)
	DeleteCoach(write http.ResponseWriter, request *http.Request)

	GetAllPupilsForCoach(write http.ResponseWriter, request *http.Request)
	CreatePupil(write http.ResponseWriter, request *http.Request)
	GetPupil(write http.ResponseWriter, request *http.Request)
	GetPupilFull(write http.ResponseWriter, request *http.Request)
	UpdatePupil(write http.ResponseWriter, request *http.Request)
	DeletePupil(write http.ResponseWriter, request *http.Request)
	GetAllPupilsList(write http.ResponseWriter, request *http.Request)

	GetCoachSchedule(write http.ResponseWriter, request *http.Request)
	CreateClass(write http.ResponseWriter, request *http.Request)

	Authorize(write http.ResponseWriter, request *http.Request)

	RecoverPupils(pupils []models.Pupil)
	RecoverCoaches(coaches []models.Coach)
	RecoverAdmins(admins []models.Admin)

	ShowCacheUsers(write http.ResponseWriter, request *http.Request)
	GetAllCoachList(write http.ResponseWriter, request *http.Request)
	GetNearestBirthdays(write http.ResponseWriter, request *http.Request)
}

type server struct {
	PupilsCache  cache.IPupils
	CoachesCache cache.ICoaches
	AdminsCache  cache.IAdmin
	Repository   *repository.Repository
}

const (
	GroupV1 = "/api/v1"

	CreateAdminRoute = "/users/admin/create"
	GetAdminRoute    = "/users/admin/get"
	GetAdminsList    = "/users/admin/list"

	CreateCoachRoute   = "/users/coach/create"
	GetCoachRoute      = "/users/coach/get"
	GetCoachFullRoute  = "/users/coach/get/full"
	UpdateCoachRoute   = "/users/coach/update"
	DeleteCoachRoute   = "/users/coach/delete"
	GetCoachPupilsList = "/users/coach/pupils"
	GetCoachesList     = "/users/coach/list"

	CreatePupilRoute  = "/users/pupil/create"
	GetPupilRoute     = "/users/pupil/get"
	GetPupilFullRoute = "/users/pupil/get/full"
	UpdatePupilRoute  = "/users/pupil/update"
	DeletePupilRoute  = "/users/pupil/delete"
	GetPupilsList     = "/users/pupil/list"

	GetCoachSchedule = "/classes/coach/schedule"
	CreateNewClass   = "/classes/new"

	AuthRoute = "/auth"

	CacheGetAllRoute = "/cache/all"

	GetPupilsBirthdayList = "/additional/birthday/list"
)

func NewServer(repo *repository.Repository) Server {
	c := cache.NewCache()
	return &server{
		PupilsCache:  c.NewPupilsCache(),
		CoachesCache: c.NewCoachesCache(),
		AdminsCache:  c.NewAdminsCache(),
		Repository:   repo,
	}
}

func (s *server) RecoverPupils(pupils []models.Pupil) {
	s.PupilsCache.WritingSession(pupils)
}

func (s *server) RecoverCoaches(coaches []models.Coach) {
	s.CoachesCache.WritingSession(coaches)
}

func (s *server) RecoverAdmins(admins []models.Admin) {
	s.AdminsCache.WritingSession(admins)
}
