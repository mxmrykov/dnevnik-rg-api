package models

import "time"

type (
	General struct {
		UDID    int
		Key     int
		DateReg time.Time
	}
	Pupil struct {
		General
	}

	Coach struct {
		General
	}

	Admin struct {
		General
	}
)
