package external

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"dnevnik-rg.ru/pkg/utils"

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
		s.Zerolog.Err(decodingBodyErr).Msg("cannot convert request body")
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	xUserIdString := request.Header.Get("X-User-Id")
	if decoded.Checksum == "" || xUserIdString == "" {
		s.Zerolog.Err(errors.New("one of required params is empty")).Str("cannot convert user id", xUserIdString)
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	userId, errConv := strconv.Atoi(xUserIdString)
	if errConv != nil {
		s.Zerolog.Err(errConv).Str("cannot convert user id", xUserIdString)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	auth, errAuthorize := s.Store.Authorize(userId, decoded.Checksum, request.RemoteAddr)
	if errAuthorize != nil {
		s.Zerolog.Err(errAuthorize).Msg("authorization error")
		write.WriteHeader(http.StatusUnauthorized)
		WriteResponse(write, "Неверные данные для входа", true, http.StatusUnauthorized)
		return
	}
	token, _ := utils.SetLongJwt(userId, decoded.Checksum, "ADMIN")
	if updErr := s.Store.AutoUpdateToken(token, userId); updErr != nil {
		s.Zerolog.Err(updErr).Msg("updating user token error")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Неверные данные для входа", true, http.StatusInternalServerError)
		return
	}

	auth.Token = token
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Данные авторизации получены", false, http.StatusOK, auth)
	return
}
