package http

import (
	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/http/external"
	"log"
	"net/http"
)

func NewHttp(configHttp *config.Http) {
	mux := http.NewServeMux()
	server := external.NewServer()
	mux.HandleFunc(external.GroupV1+external.CreateAdminRoute, server.CreateAdmin)
	handler := external.Logger(mux)
	log.Fatal(http.ListenAndServe(configHttp.Host+":"+configHttp.Port, handler))
}
