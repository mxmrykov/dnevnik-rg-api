package http

import (
	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/http/external"
	"dnevnik-rg.ru/internal/repository"
	"log"
	"net/http"
)

func NewHttp(configHttp *config.Http, repo *repository.Repository, recoveryRequired bool) {
	mux := http.NewServeMux()
	server := external.NewServer(repo)

	go func() {
		if !recoveryRequired {
			return
		}
		log.Println("recovery is required, starting...")
		pupils, errRecoveryPupils := repo.GetAllPupils()
		if errRecoveryPupils != nil {
			log.Println("error recovering pupils from DB:", errRecoveryPupils)
		}
		server.RecoverPupils(pupils)
		coaches, errRecoveryCoaches := repo.GetAllCoaches()
		if errRecoveryCoaches != nil {
			log.Println("error recovering coaches from DB:", errRecoveryCoaches)
		}
		server.RecoverCoaches(coaches)
		admins, errRecoveryAdmins := repo.GetAllAdmins()
		if errRecoveryAdmins != nil {
			log.Println("error recovering admins from DB:", errRecoveryAdmins)
		}
		server.RecoverAdmins(admins)
		log.Println("recovery is overed")
		return
	}()

	log.Println("starting web server...")

	//Group Admin
	mux.HandleFunc(external.GroupV1+external.CreateAdminRoute, server.CreateAdmin)
	mux.HandleFunc(external.GroupV1+external.GetAdminRoute, server.GetAdmin)
	mux.HandleFunc(external.GroupV1+external.GetAdminsList, server.GetAllAdminsExcept)

	//Group Coach
	mux.HandleFunc(external.GroupV1+external.CreateCoachRoute, server.CreateCoach)
	mux.Handle(external.GroupV1+external.GetCoachRoute,
		external.CheckCoachId(http.HandlerFunc(server.GetCoach)),
	)
	mux.Handle(external.GroupV1+external.GetCoachFullRoute,
		external.CheckCoachId(http.HandlerFunc(server.GetCoachFull)),
	)
	mux.Handle(external.GroupV1+external.UpdateCoachRoute,
		external.CheckCoachId(http.HandlerFunc(server.UpdateCoach)),
	)
	mux.HandleFunc(external.GroupV1+external.DeleteCoachRoute, server.DeleteCoach)
	mux.HandleFunc(external.GroupV1+external.GetCoachesList, server.GetAllCoachList)

	//Group Pupil
	mux.HandleFunc(external.GroupV1+external.CreatePupilRoute, server.CreatePupil)
	mux.Handle(external.GroupV1+external.GetPupilRoute,
		external.CheckPupilId(http.HandlerFunc(server.GetPupil)),
	)
	mux.Handle(external.GroupV1+external.GetPupilFullRoute,
		external.CheckPupilId(http.HandlerFunc(server.GetPupilFull)),
	)
	mux.Handle(external.GroupV1+external.UpdatePupilRoute,
		external.CheckPupilId(http.HandlerFunc(server.UpdatePupil)),
	)
	mux.HandleFunc(external.GroupV1+external.DeletePupilRoute, server.DeletePupil)
	mux.HandleFunc(external.GroupV1+external.GetPupilsList, server.GetAllPupilsList)

	//Group auth
	mux.HandleFunc(external.GroupV1+external.AuthRoute, server.Authorize)

	//Group Cache
	mux.HandleFunc(external.GroupV1+external.CacheGetAllRoute, server.ShowCacheUsers)

	handler := external.CheckPermission(mux)
	handler = external.Logger(handler)
	handler = external.SetCors(handler)
	log.Println("server started on", configHttp.Host+":"+configHttp.Port)
	log.Fatal(http.ListenAndServe(configHttp.Host+":"+configHttp.Port, handler))
}
