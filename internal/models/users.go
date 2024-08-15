package models

type (
	General struct {
		UDID    int    `json:"UDID"`
		Key     int    `json:"key"`
		Fio     string `json:"fio"`
		DateReg string `json:"date_reg"`
		LogoUri string `json:"logo_uri"`
		Role    string `json:"role"`
	}

	Pupil struct {
		General
		Coach        int    `json:"coach"`
		HomeCity     string `json:"home_city"`
		TrainingCity string `json:"training_city"`
		Birthday     string `json:"birthday"`
		About        string `json:"about"`
		CoachReview  string `json:"coach_review"`
	}

	Coach struct {
		General
		HomeCity     string `json:"home_city"`
		TrainingCity string `json:"training_city"`
		Birthday     string `json:"birthday"`
		About        string `json:"about"`
	}

	Admin struct {
		General
	}
)
