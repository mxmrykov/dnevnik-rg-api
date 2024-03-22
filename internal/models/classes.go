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
	Pupil         []int  `json:"pupil"`
	Coach         int    `json:"coach"`
	ClassDate     string `json:"class_date"`
	ClassTime     string `json:"class_time"`
	ClassDuration string `json:"class_duration"`
}
