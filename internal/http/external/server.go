package external

import (
	"dnevnik-rg.ru/internal/cache"
	"dnevnik-rg.ru/internal/repository"
	"net/http"
)

type Server interface {
	CreateAdmin(write http.ResponseWriter, request *http.Request)
	GetAdmin(write http.ResponseWriter, request *http.Request)
	CreateCoach(write http.ResponseWriter, request *http.Request)
	GetCoachFull(write http.ResponseWriter, request *http.Request)
	UpdateCoach(write http.ResponseWriter, request *http.Request)
	DeleteCoach(write http.ResponseWriter, request *http.Request)
}

type server struct {
	Cache      *cache.Cache
	Repository *repository.Repository
}

func NewServer(repo *repository.Repository) server {
	return server{Cache: cache.NewCache(), Repository: repo}
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
)
