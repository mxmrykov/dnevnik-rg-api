package external

import (
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		log.Printf("Handle %s | Method %s | Time %s", request.URL, request.Method, time.Since(start))
	})
}

func CheckPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		xUserId := request.Header.Get("X-User-Id")
		Auth := request.Header.Get("Authorization")
		if len(xUserId) == 0 || len(Auth) == 0 {
			log.Println("missing some required params")
			writer.WriteHeader(http.StatusBadRequest)
			WriteResponse(writer, "Ошибка авторизации", true, http.StatusBadRequest)
			return
		}
		token, errParse := jwt.Parse(Auth, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if errParse != nil {
			log.Printf("err parsing token: %v\n", errParse)
			writer.WriteHeader(http.StatusUnauthorized)
			WriteResponse(writer, "Ошибка авторизации", true, http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("token is invalid")
			if v, ok := errParse.(*jwt.ValidationError); ok {
				if v.Errors&jwt.ValidationErrorMalformed != 0 {
					log.Println("ValidationErrorMalformed")
				} else if v.Errors&jwt.ValidationErrorExpired != 0 {
					log.Println("ValidationErrorExpired")
				} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
					log.Println("ValidationErrorNotValidYet")
				}
			}
			writer.WriteHeader(http.StatusUnauthorized)
			WriteResponse(writer, "Ошибка авторизации", true, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
