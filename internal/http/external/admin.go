package external

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	if err := decoder.Decode(&decoded); err != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	key := int(utils.GetKey())
	checkSum, timeNow := utils.NewPassword(), time.Now().Format(time.RFC3339)
	byteArr := []byte(utils.HashSumGen(key, checkSum))
	cs := base64.StdEncoding.EncodeToString(byteArr)
	token, err := utils.SetLongJwt(key, cs, "ADMIN")
	if err != nil {
		s.Zerolog.Err(err).Msg("error creating new token")
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
	if err = s.Store.NewAdmin(newAdmin); err != nil {
		s.Zerolog.Err(err).Msg("error creating new admin")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	if err = s.Store.NewPassword(newPassword); err != nil {
		s.Zerolog.Err(err).Msg("error creating new password for admin")
		if err = s.Store.DeleteAdmin(key); err != nil {
			s.Zerolog.Err(err).Msg("error deleting new admin without password")
		}
		s.Zerolog.Info().Msg("new admin without password cleared")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	s.AdminsCache.WriteAdmin(newAdmin)
	admin, err := s.Store.GetAdmin(key)
	if err != nil {
		s.Zerolog.Err(err).Msg("error at returning new admin data")
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
	UserId, err := strconv.Atoi(UserIdString)
	if err != nil {
		s.Zerolog.Err(err).Str("cannot convert user id", UserIdString)
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
	admin, err := s.Store.GetAdmin(UserId)
	if err != nil {
		s.Zerolog.Err(err).Msg("cannot check admin")
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
	UserId, err := strconv.Atoi(UserIdString)
	if err != nil {
		s.Zerolog.Err(err).Str("cannot convert user id", UserIdString)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	admins, err := s.Store.GetAllAdminsExcept(UserId)
	if err != nil {
		s.Zerolog.Err(err).Msg("cannot list admins")
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

	filterCoachIDString, filterPupilIDString := request.URL.Query().Get("filterCoachID"), request.URL.Query().Get("filterPupilID")

	if ok, _ := s.checkExistence(write, request); !ok {
		return
	}

	filterCoachID, filterPupilID := new(int), new(int)

	if filterCoachIDString != "" {
		filterCoach, err := strconv.Atoi(filterCoachIDString)

		if err != nil {
			s.Zerolog.Err(err).Str("cannot convert coach id", filterCoachIDString)
			write.WriteHeader(http.StatusInternalServerError)
			WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
			return
		}

		filterCoachID = &filterCoach
	}

	if filterPupilIDString != "" {
		filterPupil, err := strconv.Atoi(filterPupilIDString)

		if err != nil {
			s.Zerolog.Err(err).Str("cannot convert pupil id", filterPupilIDString)
			write.WriteHeader(http.StatusInternalServerError)
			WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
			return
		}

		filterPupilID = &filterPupil
	}

	timeNow := time.Now()
	timeNowMonth, timeNowDay := *new(string), *new(string)

	if timeNow.Month() < 10 {
		timeNowMonth = "0" + strconv.Itoa(int(timeNow.Month()))
	} else {
		timeNowMonth = strconv.Itoa(int(timeNow.Month()))
	}

	if timeNow.Day() < 10 {
		timeNowDay = "0" + strconv.Itoa(int(timeNow.Day()))
	} else {
		timeNowDay = strconv.Itoa(int(timeNow.Day()))
	}

	date := fmt.Sprintf("%d-%s-%s", timeNow.Year(), timeNowMonth, timeNowDay)

	s.Zerolog.Info().Msg(date)

	classes, err := s.Store.GetAdminClassesHistory(date)

	if err != nil {
		s.Zerolog.Err(err).Msg("cannot get admin history")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	classesRes, err := utils.ValidateHistoryClasses(classes, filterCoachID, filterPupilID)

	if err != nil {
		s.Zerolog.Err(err).Msg("cannot validate classes history")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "История администратора получена", false, http.StatusOK, classesRes)
	return
}
