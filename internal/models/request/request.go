package requests

import "github.com/golang-jwt/jwt"

type (
	NewAdmin struct {
		Fio string `json:"fio"`
	}
	JwtPayload struct {
		Key      int    `json:"key"`
		CheckSum string `json:"check_sum"`
		Role     string `json:"role"`
		jwt.StandardClaims
	}
)
