package external

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/pkg/utils"
)

func (s *server) CreateAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	decoder := json.NewDecoder(request.Body)
	var decoded requests.NewAdmin
	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	key := int(utils.GetKey())
	checkSum := utils.NewPassword()
	timeNow := time.Now().Format(time.RFC3339)
	byteArr := []byte(utils.HashSumGen(key, checkSum))
	cs := base64.StdEncoding.EncodeToString(byteArr)
	token, errCreateToken := utils.SetLongJwt(key, cs, "ADMIN")
	if errCreateToken != nil {
		s.Zerolog.Err(errCreateToken).Msg("error creating new token")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	newAdmin := models.Admin{
		General: models.General{
			Key:     key,
			Fio:     decoded.Fio,
			DateReg: timeNow,
			LogoUri: "https://dnevnik-rg.ru/admin-logo.png",
			Role:    "ADMIN",
		},
	}
	newPassword := models.Password{
		Key:        key,
		CheckSum:   cs,
		LastUpdate: timeNow,
		Token:      token,
	}
	if errNewAdmin := s.Store.NewAdmin(newAdmin); errNewAdmin != nil {
		s.Zerolog.Err(errNewAdmin).Msg("error creating new admin")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	if errNewPassword := s.Store.NewPassword(newPassword); errNewPassword != nil {
		s.Zerolog.Err(errNewPassword).Msg("error creating new password for admin")
		if errClearingBrokenAdmin := s.Store.DeleteAdmin(key); errClearingBrokenAdmin != nil {
			s.Zerolog.Err(errClearingBrokenAdmin).Msg("error deleting new admin without password")
		}
		s.Zerolog.Info().Msg("new admin without password cleared")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	s.AdminsCache.WriteAdmin(newAdmin)
	admin, errGetAdmin := s.Store.GetAdmin(key)
	if errGetAdmin != nil {
		s.Zerolog.Err(errGetAdmin).Msg("error at returning new admin data")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Не удалось получить созданного администратора", true, http.StatusInternalServerError)
		return
	}
	admin.Private.CheckSum = checkSum
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Администратор зарегистрирован", false, http.StatusOK, admin)
	return
}

func (s *server) GetAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	UserIdString := request.Header.Get("X-User-Id")
	UserId, errConv := strconv.Atoi(UserIdString)
	if errConv != nil {
		s.Zerolog.Err(errConv).Str("cannot convert user id", UserIdString)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	if p, ok_ := s.AdminsCache.ReadById(UserId); ok_ {
		s.Zerolog.Info().Int("admin loaded from cache", (*p).Key)
		write.WriteHeader(http.StatusOK)
		WriteDataResponse(write, "Администратор получен", false, http.StatusOK, *p)
		return
	}
	admin, errGetAdmin := s.Store.GetAdmin(UserId)
	if errGetAdmin != nil {
		s.Zerolog.Err(errGetAdmin).Msg("cannot check admin")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка сервера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Администратор получен", false, http.StatusOK, admin)
	return
}

func (s *server) GetAllAdminsExcept(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	UserIdString := request.Header.Get("X-User-Id")
	UserId, errConv := strconv.Atoi(UserIdString)
	if errConv != nil {
		s.Zerolog.Err(errConv).Str("cannot convert user id", UserIdString)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	admins, errGetAdmin := s.Store.GetAllAdminsExcept(UserId)
	if errGetAdmin != nil {
		s.Zerolog.Err(errGetAdmin).Msg("cannot list admins")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка сервера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список администраторов получен", false, http.StatusOK, admins)
	return
}

func (s *server) GetClassesHistoryForAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}

}
