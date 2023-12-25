package models

type Password struct {
	UDID       int    `json:"udid"`
	Key        int    `json:"key"`
	CheckSum   string `json:"check_sum"`
	LastUpdate string `json:"last_update"`
	Token      string `json:"token"`
}
