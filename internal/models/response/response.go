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
		models.General
		Private Private `json:"private"`
	}
	AdminList struct {
		Key     int    `json:"key"`
		Fio     string `json:"fio"`
		LogoUri string `json:"logo_uri"`
	}
	CoachList struct {
		Key     int    `json:"key"`
		Fio     string `json:"fio"`
		LogoUri string `json:"logo_uri"`
	}
	PupilList struct {
		Key     int    `json:"key"`
		Fio     string `json:"fio"`
		LogoUri string `json:"logo_uri"`
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

	CanceledClass struct {
		Canceled bool `json:"canceled"`
		Key      int  `json:"key"`
	}

	DeletedClass struct {
		Deleted bool `json:"deleted"`
		Key     int  `json:"key"`
	}
)
