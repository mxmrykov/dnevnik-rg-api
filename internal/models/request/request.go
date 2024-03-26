package requests

import "github.com/golang-jwt/jwt"

type (
	NewAdmin struct {
		Fio string `json:"fio"`
	}
	NewCoach struct {
		Fio          string `json:"fio"`
		HomeCity     string `json:"home_city"`
		TrainingCity string `json:"training_city"`
		Birthday     string `json:"birthday"`
		About        string `json:"about"`
	}
	NewPupil struct {
		Fio          string `json:"fio"`
		HomeCity     string `json:"home_city"`
		TrainingCity string `json:"training_city"`
		Birthday     string `json:"birthday"`
		About        string `json:"about"`
		Coach        int    `json:"coach"`
		CoachReview  string `json:"coach_review"`
	}
	JwtPayload struct {
		Key      int    `json:"key"`
		CheckSum string `json:"check_sum"`
		Role     string `json:"role"`
		jwt.StandardClaims
	}
	UpdateCoach struct {
		Fio          string `json:"fio"`
		HomeCity     string `json:"home_city"`
		TrainingCity string `json:"training_city"`
		Birthday     string `json:"birthday"`
		About        string `json:"about"`
	}
	UpdatePupil struct {
		Fio          string `json:"fio"`
		HomeCity     string `json:"home_city"`
		TrainingCity string `json:"training_city"`
		Birthday     string `json:"birthday"`
		About        string `json:"about"`
		CoachReview  string `json:"coach_review"`
	}
	Auth struct {
		Checksum string `json:"checksum"`
	}
	BirthDayList struct {
		Key      int    `json:"key"`
		Fio      string `json:"fio"`
		Birthday string `json:"birthday"`
	}
	CreateClass struct {
		Pupil     []int  `json:"pupil"`
		Capacity  int    `json:"capacity"`
		Coach     int    `json:"coach"`
		ClassDate string `json:"classDate"`
		ClassTime string `json:"classTime"`
		Duration  string `json:"duration"`
		Price     int    `json:"price"`
		ClassType string `json:"classType"`
		IsOpen    bool   `json:"isOpen"`
	}
)
