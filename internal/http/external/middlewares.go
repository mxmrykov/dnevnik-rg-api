package external

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		log.Printf("Handle %s | Method %s | Reqponse %d", request.URL, request.Method, request.Response.StatusCode)
		next.ServeHTTP(write, request)
	})
}
