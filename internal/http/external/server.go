package external

import (
	"dnevnik-rg.ru/internal/cache"
	"dnevnik-rg.ru/internal/models"
	"dnevnik-rg.ru/internal/repository"
	"net/http"
)

type Server interface {
	CreateAdmin(write http.ResponseWriter, request *http.Request)
	GetAdmin(write http.ResponseWriter, request *http.Request)
	CreateCoach(write http.ResponseWriter, request *http.Request)
	GetCoach(write http.ResponseWriter, request *http.Request)
	GetCoachFull(write http.ResponseWriter, request *http.Request)
	UpdateCoach(write http.ResponseWriter, request *http.Request)
	DeleteCoach(write http.ResponseWriter, request *http.Request)
	CreatePupil(write http.ResponseWriter, request *http.Request)
	GetPupil(write http.ResponseWriter, request *http.Request)
	GetPupilFull(write http.ResponseWriter, request *http.Request)
	UpdatePupil(write http.ResponseWriter, request *http.Request)
	DeletePupil(write http.ResponseWriter, request *http.Request)
	RecoverPupils(pupils []models.Pupil)
	RecoverCoaches(coaches []models.Coach)
	RecoverAdmins(admins []models.Admin)
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

	CreateCoachRoute  = "/users/coach/create"
	GetCoachRoute     = "/users/coach/get"
	GetCoachFullRoute = "/users/coach/get/full"
	UpdateCoachRoute  = "/users/coach/update"
	DeleteCoachRoute  = "/users/coach/delete"

	CreatePupilRoute  = "/users/pupil/create"
	GetPupilRoute     = "/users/pupil/get"
	GetPupilFullRoute = "/users/pupil/get/full"
	UpdatePupilRoute  = "/users/pupil/update"
	DeletePupilRoute  = "/users/pupil/delete"
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