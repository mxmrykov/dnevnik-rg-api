package http

import (
	"net/http"
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/http/external"
	"dnevnik-rg.ru/internal/store"
	"github.com/rs/zerolog"
)

func NewHttp(configHttp *config.Http, rgStore store.Store, recoveryRequired bool, logger *zerolog.Logger) {
	mux := http.NewServeMux()
	server := external.NewServer(rgStore, logger)
	cacheRecoveryTimeStart := time.Now()

	go func(cacheRecoveryTimeStart *time.Time) {
		if !recoveryRequired {
			return
		}
		logger.Info().Msg("cache recovery is required, starting...")
		pupils, errRecoveryPupils := rgStore.GetAllPupils()
		if errRecoveryPupils != nil {
			logger.Err(errRecoveryPupils).Msg("error recovering pupils from DB to cache")
		}
		server.RecoverPupils(pupils)
		coaches, errRecoveryCoaches := rgStore.GetAllCoaches()
		if errRecoveryCoaches != nil {
			logger.Err(errRecoveryCoaches).Msg("error recovering coaches from DB to cache")
		}
		server.RecoverCoaches(coaches)
		admins, errRecoveryAdmins := rgStore.GetAllAdmins()
		if errRecoveryAdmins != nil {
			logger.Err(errRecoveryCoaches).Msg("error recovering admins from DB to cache")
		}
		server.RecoverAdmins(admins)
		logger.Info().Dur("recovery time", time.Since(*cacheRecoveryTimeStart)).Msg("cache recovery is overed")
		return
	}(&cacheRecoveryTimeStart)

	logger.Info().Msg("starting web server")

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
	mux.HandleFunc(external.GroupV1+external.GetArchiveCoaches, server.ArchiveCoachGet)
	mux.HandleFunc(external.GroupV1+external.ArchiveCoachRoute, server.ArchiveCoach)
	mux.Handle(external.GroupV1+external.DeArchiveCoach,
		external.CheckCoachId(http.HandlerFunc(server.DearchiveCoach)),
	)

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
	mux.HandleFunc(external.GroupV1+external.ArchivePupilRoute, server.ArchivePupil)
	mux.HandleFunc(external.GroupV1+external.GetArchivePupils, server.ArchivePupilGet)
	mux.HandleFunc(external.GroupV1+external.GetPupilsList, server.GetAllPupilsList)
	mux.Handle(external.GroupV1+external.DeArchivePupil,
		external.CheckPupilId(http.HandlerFunc(server.DearchivePupil)),
	)

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
	mux.HandleFunc(external.GroupV1+external.ClassesHistoryAdmin, server.GetClassesHistoryForAdmin)

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

	logger.Err(http.ListenAndServe(configHttp.Host+":"+configHttp.Port, handler)).Dur("server work time", time.Since(cacheRecoveryTimeStart))
}
