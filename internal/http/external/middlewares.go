package external

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(write, request)
		log.Printf("Handle %s | Method %s | Time %s", request.URL, request.Method, time.Since(start))
	})
}
