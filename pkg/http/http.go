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
	//Group Admin
	mux.HandleFunc(external.GroupV1+external.CreateAdminRoute, server.CreateAdmin)
	mux.HandleFunc(external.GroupV1+external.GetAdminRoute, server.GetAdmin)
	//Group Coach
	mux.HandleFunc(external.GroupV1+external.CreateCoachRoute, server.CreateCoach)
	mux.Handle(external.GroupV1+external.GetCoachRoute,
		external.CheckCoachId(http.HandlerFunc(server.GetCoach)),
	)
	mux.Handle(external.GroupV1+external.GetCoachFullRoute,
		external.CheckCoachId(http.HandlerFunc(server.GetCoachFull)),
	)
	mux.HandleFunc(external.GroupV1+external.UpdateCoachRoute, server.UpdateCoach)
	mux.HandleFunc(external.GroupV1+external.DeleteCoachRoute, server.DeleteCoach)
	handler := external.CheckPermission(mux)
	handler = external.Logger(handler)
	log.Println("server started on", configHttp.Host+":"+configHttp.Port)
	log.Fatal(http.ListenAndServe(configHttp.Host+":"+configHttp.Port, handler))
}
