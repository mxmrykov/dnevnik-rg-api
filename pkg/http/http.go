package http

import (
	"log"
	"net/http"
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/http/external"
	"dnevnik-rg.ru/internal/store"
)

func NewHttp(configHttp *config.Http, rgStore store.Store, recoveryRequired bool) {
	mux := http.NewServeMux()
	server := external.NewServer(rgStore)
	cacheRecoveryTimeStart := time.Now()

	go func(cacheRecoveryTimeStart *time.Time) {
		if !recoveryRequired {
			return
		}
		log.Println("cache recovery is required, starting...")
		pupils, errRecoveryPupils := rgStore.GetAllPupils()
		if errRecoveryPupils != nil {
			log.Println("error recovering pupils from DB to cache:", errRecoveryPupils)
		}
		server.RecoverPupils(pupils)
		coaches, errRecoveryCoaches := rgStore.GetAllCoaches()
		if errRecoveryCoaches != nil {
			log.Println("error recovering coaches from DB to cache:", errRecoveryCoaches)
		}
		server.RecoverCoaches(coaches)
		admins, errRecoveryAdmins := rgStore.GetAllAdmins()
		if errRecoveryAdmins != nil {
			log.Println("error recovering admins from DB to cache:", errRecoveryAdmins)
		}
		server.RecoverAdmins(admins)
		log.Println("cache recovery is overed | recovery time:", time.Since(*cacheRecoveryTimeStart))
		return
	}(&cacheRecoveryTimeStart)

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
	mux.Handle(external.GroupV1+external.GetCoachPupilsList,
		external.CheckCoachId(http.HandlerFunc(server.GetAllPupilsForCoach)),
	)
	mux.HandleFunc(external.GroupV1+external.GetCoachesList, server.GetAllCoachList)

	//Group Pupils
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

	//Group Classes
	mux.HandleFunc(external.GroupV1+external.GetCoachSchedule, server.GetCoachSchedule)
	mux.HandleFunc(external.GroupV1+external.CreateNewClass, server.CreateClass)
	mux.HandleFunc(external.GroupV1+external.GetClassesForDayAdmin, server.GetClassesTodayAdmin)
	mux.HandleFunc(external.GroupV1+external.GetClassesForDayCoach, server.GetClassesTodayCoach)
	mux.HandleFunc(external.GroupV1+external.GetClassesForDayPupil, server.GetClassesTodayPupil)
	mux.HandleFunc(external.GroupV1+external.GetClassesForMonthAdmin, server.GetClassesMonthAdmin)
	mux.HandleFunc(external.GroupV1+external.GetClassesForMonthCoach, server.GetClassesMonthCoach)
	mux.HandleFunc(external.GroupV1+external.GetClassesForMonthPupil, server.GetClassesMonthPupil)
	mux.HandleFunc(external.GroupV1+external.CancelClass, server.CancelClass)
	mux.HandleFunc(external.GroupV1+external.DeleteClass, server.DeleteClass)
	mux.HandleFunc(external.GroupV1+external.ClassInfoAdmin, server.GetClassInfoAdmin)

	//Group auth
	mux.HandleFunc(external.GroupV1+external.AuthRoute, server.Authorize)

	//Group Cache
	mux.HandleFunc(external.GroupV1+external.CacheGetAllRoute, server.ShowCacheUsers)

	//Group Additional
	mux.Handle(external.GroupV1+external.GetPupilsBirthdayList,
		external.CheckCoachId(http.HandlerFunc(server.GetNearestBirthdays)),
	)

	handler := external.CheckPermission(mux)
	handler = external.Logger(handler)
	handler = external.SetCors(handler)
	log.Fatal(http.ListenAndServe(configHttp.Host+":"+configHttp.Port, handler), "server work time:", time.Since(cacheRecoveryTimeStart))
}
