package http

import (
	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/http/external"
	"dnevnik-rg.ru/internal/repository"
	"log"
	"net/http"
)

func NewHttp(configHttp *config.Http, repo *repository.Repository) {
	mux := http.NewServeMux()
	server := external.NewServer(repo)
	mux.HandleFunc(external.GroupV1+external.CreateAdminRoute, server.CreateAdmin)
	handler := external.Logger(mux)
	log.Println("server started on", configHttp.Host+":"+configHttp.Port)
	log.Fatal(http.ListenAndServe(configHttp.Host+":"+configHttp.Port, handler))
}
