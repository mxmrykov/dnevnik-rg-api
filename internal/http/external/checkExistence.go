package external

import (
	"log"
	"net/http"
	"strconv"
)

func (s *server) checkExistence(write http.ResponseWriter, request *http.Request) (bool, error) {
	UserIdString := request.Header.Get("X-User-Id")
	UserId, errConv := strconv.Atoi(UserIdString)
	if errConv != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return false, errConv
	}
	if ok, errCheckAdmin := s.Repository.IsAdminExists(UserId); !ok || errCheckAdmin != nil {
		log.Printf("err check admin: %v\n", errCheckAdmin)
		write.WriteHeader(http.StatusForbidden)
		WriteResponse(write, "Доступ запрещен", true, http.StatusForbidden)
		return false, errCheckAdmin
	}
	return true, nil
}
