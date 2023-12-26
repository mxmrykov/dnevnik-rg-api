package app

import (
	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/repository"
	"dnevnik-rg.ru/pkg/http"
	"dnevnik-rg.ru/pkg/postgres"
	"log"
)

func App(appConfig *config.Config) {
	postgresConnection, errConnectPostgres := postgres.NewPostgres(&appConfig.Postgres)
	if errConnectPostgres != nil {
		log.Fatalf("cannot connect to postrges: %v", errConnectPostgres)
		return
	}
	repos := repository.NewRepository(postgresConnection)
	if errInitPupils := repos.InitTablePupils(); errInitPupils != nil {
		log.Printf("error initializing pupils table: %v\n", errInitPupils)
	}
	if errInitCoaches := repos.InitTableCoaches(); errInitCoaches != nil {
		log.Printf("error initializing coaches table: %v\n", errInitCoaches)
	}
	if errInitPasswords := repos.InitTablePasswords(); errInitPasswords != nil {
		log.Printf("error initializing passwords table: %v\n", errInitPasswords)
	}
	if errInitClasses := repos.InitTableClasses(); errInitClasses != nil {
		log.Printf("error initializing classes table: %v\n", errInitClasses)
	}
	if errInitAdmins := repos.InitTableAdmins(); errInitAdmins != nil {
		log.Printf("error initializing classes table: %v\n", errInitAdmins)
	}
	log.Println("db tables initialized")
	http.NewHttp(&appConfig.Http, repos, true)
}
