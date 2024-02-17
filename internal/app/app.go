package app

import (
	"context"
	"log"
	"strconv"
	"time"

	"dnevnik-rg.ru/config"
	"dnevnik-rg.ru/internal/repository"
	"dnevnik-rg.ru/pkg/http"
	"dnevnik-rg.ru/pkg/postgres"
	"dnevnik-rg.ru/pkg/utils"
	"github.com/jackc/pgx/v4/pgxpool"
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
	go pingShards(&appConfig.TgBots, postgresConnection)
	http.NewHttp(&appConfig.Http, repos, true)
}

func pingShards(botsConfig *config.TgBots, shardsConnections []*pgxpool.Pool) {
	log.Println("starting postgres pinger")
	var (
		attemptsFirst  = 0
		attemptsSecond = 0
		attemptsThird  = 0
		timeSincePing  time.Time
		err            error
	)

	for {
		timeSincePing = time.Now()
		err = shardsConnections[0].Ping(context.Background())

		if err != nil {
			attemptsFirst += 1
			log.Println("Failed to ping first shard! Attempt: ", attemptsFirst, "|", "Pinging time: ", time.Since(timeSincePing))
			switch attemptsFirst {
			case 5, 15, 30, 50, 100:
				log.Println("sending alert to telegram bot...")
				utils.SendTgTechPgPingAlert(
					botsConfig.TgTechBot.Token,
					"pg-shard-1",
					strconv.Itoa(attemptsFirst),
				)
			}
		}

		timeSincePing = time.Now()
		err = shardsConnections[1].Ping(context.Background())

		if err != nil {
			attemptsSecond += 1
			log.Println("Failed to ping first shard! Attempt: ", attemptsSecond, "|", "Pinging time: ", time.Since(timeSincePing))
			switch attemptsSecond {
			case 5, 15, 30, 50, 100:
				log.Println("sending alert to telegram bot...")
				utils.SendTgTechPgPingAlert(
					botsConfig.TgTechBot.Token,
					"pg-shard-2",
					strconv.Itoa(attemptsSecond),
				)
			}
		}

		timeSincePing = time.Now()
		err = shardsConnections[2].Ping(context.Background())

		if err != nil {
			attemptsThird += 1
			log.Println("Failed to ping first shard! Attempt: ", attemptsThird, "|", "Pinging time: ", time.Since(timeSincePing))
			switch attemptsThird {
			case 5, 15, 30, 50, 100:
				log.Println("sending alert to telegram bot...")
				utils.SendTgTechPgPingAlert(
					botsConfig.TgTechBot.Token,
					"pg-shard-3",
					strconv.Itoa(attemptsThird),
				)
			}
		}

		<-time.After(5 * time.Second)
	}
}

func serveTgTechBot() {

}
