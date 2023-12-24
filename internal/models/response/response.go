package response

type (
	Response struct {
		Data       interface{} `json:"data"`
		StatusCode int         `json:"status_code"`
		Message    string      `json:"message"`
		IsError    bool        `json:"isError"`
	}
	Admin struct {
		Key        int    `json:"key"`
		Fio        string `json:"fio"`
		DateReg    string `json:"date_reg"`
		LogoUri    string `json:"logo_uri"`
		CheckSum   string `json:"checksum"`
		LastUpdate string `json:"last_update"`
	}
)
