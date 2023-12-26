package external

import (
	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
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
	token, errCreateToken := utils.SetLongJwt(key, checkSum, timeNow)
	if errCreateToken != nil {
		log.Printf("error creating new token: %v\n", errCreateToken)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	newAdmin := models.Admin{
		General: models.General{Key: key,
			Fio:     decoded.Fio,
			DateReg: timeNow,
			LogoUri: "https://dnevnik-rg.ru/admin-logo.png",
			Role:    "ADMIN"},
	}
	newPassword := models.Password{
		Key:        key,
		CheckSum:   checkSum,
		LastUpdate: timeNow,
		Token:      token,
	}
	if errNewAdmin := s.Repository.NewAdmin(newAdmin); errNewAdmin != nil {
		log.Printf("error creating new admin: %v\n", errNewAdmin)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	if errNewPassword := s.Repository.NewPassword(newPassword); errNewPassword != nil {
		log.Printf("error creating new password for admin: %v\n", errNewPassword)
		if errClearingBrokenAdmin := s.Repository.DeleteAdmin(key); errClearingBrokenAdmin != nil {
			log.Printf("error deleting new admin without password: %v\n", errClearingBrokenAdmin)
		}
		log.Println("new admin without password cleared")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	admin, errGetAdmin := s.Repository.GetAdmin(key)
	if errGetAdmin != nil {
		log.Printf("error returns new admin data: %v\n", errGetAdmin)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Не удалось получить созданного администратора", true, http.StatusInternalServerError)
		return
	}
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
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	if p, ok_ := s.PupilsCache.ReadById(UserId); ok_ {
		log.Printf("admin loaded from cache: %d", (*p).Key)
		write.WriteHeader(http.StatusOK)
		WriteDataResponse(write, "Ученица получена", false, http.StatusOK, *p)
		return
	}
	admin, errGetAdmin := s.Repository.GetAdmin(UserId)
	if errGetAdmin != nil {
		log.Printf("cannot check admin: %v\n", errGetAdmin)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка сервера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Администратор получен", false, http.StatusOK, admin)
	return
}
