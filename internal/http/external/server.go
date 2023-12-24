package external

import (
	"dnevnik-rg.ru/internal/cache"
	"net/http"
)

type Server interface {
	CreateAdmin(write http.ResponseWriter, request *http.Request)
}

type server struct {
	Cache *cache.Cache
}

func NewServer() server {
	return server{Cache: cache.NewCache()}
}

const (
	GroupV1 = "/api/v1"

	CreateAdminRoute = "/users/admin/create"
)
