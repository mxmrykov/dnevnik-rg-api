package models

type Class struct {
	ClassMainInfo
	Presence  bool   `json:"presence"`
	Price     int    `json:"price"`
	Mark      int    `json:"mark"`
	Review    string `json:"review"`
	Scheduled bool   `json:"scheduled"`
}

type ClassMainInfo struct {
	Key           int    `json:"key"`
	Pupils        []int  `json:"pupil"`
	Coach         int    `json:"coach"`
	ClassDate     string `json:"class_date"`
	ClassTime     string `json:"class_time"`
	ClassDuration string `json:"class_duration"`
}

type ShortClassInfo struct {
	Key            int    `json:"key"`
	Pupils         []int  `json:"pupil"`
	Coach          int    `json:"coach"`
	ClassTime      string `json:"class_time"`
	ClassDuration  string `json:"class_duration"`
	ClassType      string `json:"class_type"`
	PupilCount     int    `json:"pupil_count"`
	Scheduled      bool   `json:"scheduled"`
	IsOpenToSignUp bool   `json:"is_open_to_sign_up"`
}

type ShortStringClassInfo struct {
	Key            int      `json:"key"`
	Pupils         []string `json:"pupil"`
	Coach          string   `json:"coach"`
	ClassTime      string   `json:"class_time"`
	ClassDuration  string   `json:"class_duration"`
	ClassType      string   `json:"class_type"`
	PupilCount     int      `json:"pupil_count"`
	Scheduled      bool     `json:"scheduled"`
	IsOpenToSignUp bool     `json:"is_open_to_sign_up"`
}
