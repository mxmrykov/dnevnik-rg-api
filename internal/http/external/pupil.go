package external

import (
	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (s *server) CreatePupil(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	ok, _ := s.checkCoachExistence(write, request)
	if !ok {
		return
	}
	decoder := json.NewDecoder(request.Body)
	var decoded requests.NewPupil
	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	UserIdString := request.Header.Get("X-User-Id")
	UserId, errConv := strconv.Atoi(UserIdString)
	if errConv != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	key := int(utils.GetKey())
	checkSum := utils.NewPassword()
	timeNow := time.Now().Format(time.RFC3339)
	token, errCreateToken := utils.SetLongJwt(key, checkSum, timeNow)
	bday, errBday := time.Parse("2006-01-02", decoded.Birthday)
	if errBday != nil {
		log.Printf("invalid bday format: %v\n", errBday)
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный формат данных", true, http.StatusBadRequest)
		return
	}
	if errCreateToken != nil {
		log.Printf("error creating new token: %v\n", errCreateToken)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
		return
	}
	newPupil := models.Pupil{
		General: models.General{Key: key,
			Fio:     decoded.Fio,
			DateReg: timeNow,
			LogoUri: "https://dnevnik-rg.ru/pupil-logo.png",
		},
		Coach:        UserId,
		HomeCity:     decoded.HomeCity,
		TrainingCity: decoded.TrainingCity,
		Birthday:     bday.Format(time.RFC3339),
		About:        decoded.About,
		CoachReview:  decoded.CoachReview,
	}
	newPassword := models.Password{
		Key:        key,
		CheckSum:   checkSum,
		LastUpdate: timeNow,
		Token:      token,
	}
	if errNewCoach := s.Repository.CreatePupil(newPupil); errNewCoach != nil {
		log.Printf("error creating new pupil: %v\n", errNewCoach)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания ученицы", true, http.StatusInternalServerError)
		return
	}
	if errNewPassword := s.Repository.NewPassword(newPassword); errNewPassword != nil {
		log.Printf("error creating new password for coach: %v\n", errNewPassword)
		if errClearingBrokenAdmin := s.Repository.DeleteAdmin(key); errClearingBrokenAdmin != nil {
			log.Printf("error deleting new pupil without password: %v\n", errClearingBrokenAdmin)
		}
		log.Println("new pupil without password cleared")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания ученицы", true, http.StatusInternalServerError)
		return
	}
	pupil, errGetCoach := s.Repository.GetPupilFull(key)
	if errGetCoach != nil {
		log.Printf("error returns new pupil data: %v\n", errGetCoach)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Не удалось получить созданную ученицу", true, http.StatusInternalServerError)
		return
	}
	s.PupilsCache.WritePupil(newPupil)
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Администратор зарегистрирован", false, http.StatusOK, pupil)
	return
}

func (s *server) GetPupil(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	pupilIdString, userId := request.URL.Query().Get("pupilId"), request.Header.Get("X-User-Id")
	if !strings.EqualFold(pupilIdString, userId) {
		ok, _ := s.checkCoachExistence(write, request)
		if !ok {
			return
		}
	}
	pupilId, errConvCoach := strconv.Atoi(pupilIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	if p, ok := s.PupilsCache.ReadById(pupilId); ok {
		log.Printf("pupil loaded from cache: %d", (*p).Key)
		write.WriteHeader(http.StatusOK)
		WriteDataResponse(write, "Ученица получена", false, http.StatusOK, *p)
		return
	}
	pupil, errGetPupil := s.Repository.GetPupil(pupilId)
	if errGetPupil != nil {
		log.Printf("error returns new pupil data: %v\n", errGetPupil)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось получить ученицу", true, http.StatusNotFound)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Ученица получена", false, http.StatusOK, pupil)
	return
}

func (s *server) GetPupilFull(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	pupilIdString, userId := request.URL.Query().Get("pupilId"), request.Header.Get("X-User-Id")
	if !strings.EqualFold(pupilIdString, userId) {
		ok, _ := s.checkCoachExistence(write, request)
		if !ok {
			return
		}
	}
	pupilId, errConvCoach := strconv.Atoi(pupilIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	pupil, errGetPupil := s.Repository.GetPupilFull(pupilId)
	if errGetPupil != nil {
		log.Printf("error returns new pupil data: %v\n", errGetPupil)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось получить ученицу", true, http.StatusNotFound)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Ученица получена", false, http.StatusOK, pupil)
	return
}

func (s *server) UpdatePupil(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	pupilId, userId := request.URL.Query().Get("pupilId"), request.Header.Get("X-User-Id")
	if !strings.EqualFold(pupilId, userId) {
		ok, _ := s.checkCoachExistence(write, request)
		if !ok {
			return
		}
	}
	decoder := json.NewDecoder(request.Body)
	var decoded requests.UpdatePupil
	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	if decoded.Birthday != "" {
		bday, errBday := time.Parse("2006-01-02", decoded.Birthday)
		if errBday != nil {
			log.Printf("invalid bday format: %v\n", errBday)
			write.WriteHeader(http.StatusBadRequest)
			WriteResponse(write, "Неверный формат данных", true, http.StatusBadRequest)
			return
		}
		decoded.Birthday = bday.Format(time.RFC3339)
	}
	reflectBody := reflect.ValueOf(decoded)
	var (
		params []string
		values []string
	)
	for i := 0; i < reflectBody.NumField(); i += 1 {
		if reflectBody.Field(i).Interface() != "" {
			key := reflectBody.Type().Field(i).Tag.Get("json")
			keyBaseName := reflectBody.Type().Field(i).Name
			value := reflectBody.FieldByName(keyBaseName).Interface().(string)
			params = append(params, key)
			values = append(values, value)
		}
	}
	if len(params) == 0 || len(params) != len(values) {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Не указаны данные для обновления", true, http.StatusBadRequest)
		return
	}
	pupilIdInt, errConvPupil := strconv.Atoi(pupilId)
	if errConvPupil != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	sql := utils.GenerateUpdateSql("pupil", pupilIdInt, params, values)
	errUpdatePupil := s.Repository.UpdatePupil(sql)
	if errUpdatePupil != nil {
		log.Printf("error returns new coach data: %v\n", errUpdatePupil)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка обновления ученицы", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteResponse(write, "Ученица обновлена", false, http.StatusOK)
	return
}

func (s *server) DeletePupil(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	pupilIdString := request.URL.Query().Get("pupilId")
	ok, _ := s.checkCoachExistence(write, request)
	if !ok {
		return
	}
	pupilId, errConvCoach := strconv.Atoi(pupilIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	errDeletePupil := s.Repository.DeletePupil(pupilId)
	if errDeletePupil != nil {
		log.Printf("error returns delete pupil: %v\n", errDeletePupil)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось удалить ученицу", true, http.StatusNotFound)
		return
	}
	s.PupilsCache.RemovePupil(pupilId)
	write.WriteHeader(http.StatusOK)
	WriteResponse(write, "Ученица удалена", false, http.StatusOK)
	return
}