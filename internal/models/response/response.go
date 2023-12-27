package response

import (
	"dnevnik-rg.ru/internal/models"
)

type (
	Response struct {
		Data       interface{} `json:"data"`
		StatusCode int         `json:"status_code"`
		Message    string      `json:"message"`
		IsError    bool        `json:"isError"`
	}
	Private struct {
		CheckSum   string `json:"checksum"`
		LastUpdate string `json:"last_update"`
		Token      string `json:"token"`
	}
	Admin struct {
		Key     int     `json:"key"`
		Fio     string  `json:"fio"`
		DateReg string  `json:"date_reg"`
		LogoUri string  `json:"logo_uri"`
		Private Private `json:"private"`
	}
	CoachFull struct {
		models.Coach
		Private Private `json:"private"`
	}
	Coach struct {
		models.Coach
	}
	PupilFull struct {
		models.Pupil
		Private Private `json:"private"`
	}
	Pupil struct {
		models.Pupil
	}
	Auth struct {
		Key   int    `json:"key"`
		Token string `json:"token"`
		Role  string `json:"role"`
	}
)
