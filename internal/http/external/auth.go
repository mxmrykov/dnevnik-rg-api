package external

import (
	"dnevnik-rg.ru/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	requests "dnevnik-rg.ru/internal/models/request"
)

func (s *server) Authorize(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(request.Body)
	var decoded requests.Auth
	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	xUserIdString := request.Header.Get("X-User-Id")
	if decoded.Checksum == "" || xUserIdString == "" {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	userId, errConv := strconv.Atoi(xUserIdString)
	if errConv != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	auth, errAuthorize := s.Repository.Authorize(userId, decoded.Checksum)
	if errAuthorize != nil {
		log.Printf("error returns auth user: %v\n", errAuthorize)
		write.WriteHeader(http.StatusUnauthorized)
		WriteResponse(write, "Неверные данные для входа", true, http.StatusUnauthorized)
		return
	}
	token, _ := utils.SetLongJwt(userId, decoded.Checksum, "ADMIN")
	if updErr := s.Repository.AutoUpdateToken(token, userId); updErr != nil {
		log.Printf("error returns upd user token: %v\n", updErr)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Неверные данные для входа", true, http.StatusInternalServerError)
		return
	}
	auth, _ = s.Repository.Authorize(userId, decoded.Checksum)
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Данные авторизации получены", false, http.StatusOK, auth)
	return
}
