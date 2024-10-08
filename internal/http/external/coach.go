package external

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/pkg/utils"
)

func (s *server) CreateCoach(write http.ResponseWriter, request *http.Request) {
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
	var decoded requests.NewCoach
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
	token, errCreateToken := utils.SetLongJwt(key, cs, "COACH")
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
	newCoach := models.Coach{
		General: models.General{Key: key,
			Fio:     decoded.Fio,
			DateReg: timeNow,
			LogoUri: "https://dnevnik-rg.ru/coach-logo.png",
		},
		HomeCity:     decoded.HomeCity,
		TrainingCity: decoded.TrainingCity,
		Birthday:     bday.Format(time.RFC3339),
		About:        decoded.About,
	}
	newPassword := models.Password{
		Key:        key,
		CheckSum:   cs,
		LastUpdate: timeNow,
		Token:      token,
	}
	if errNewCoach := s.Store.CreateCoach(newCoach); errNewCoach != nil {
		log.Printf("error creating new coach: %v\n", errNewCoach)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания тренера", true, http.StatusInternalServerError)
		return
	}
	if errNewPassword := s.Store.NewPassword(newPassword); errNewPassword != nil {
		log.Printf("error creating new password for coach: %v\n", errNewPassword)
		if errClearingBrokenAdmin := s.Store.DeleteAdmin(key); errClearingBrokenAdmin != nil {
			log.Printf("error deleting new coach without password: %v\n", errClearingBrokenAdmin)
		}
		log.Println("new coach without password cleared")
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка создания тренера", true, http.StatusInternalServerError)
		return
	}
	coach, errGetCoach := s.Store.GetCoachFull(key)
	if errGetCoach != nil {
		log.Printf("error returns new coach data: %v\n", errGetCoach)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Не удалось получить созданного тренера", true, http.StatusInternalServerError)
		return
	}
	coach.Private.CheckSum = checkSum
	s.CoachesCache.WriteCoach(newCoach)
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Тренер зарегистрирован", false, http.StatusOK, coach)
	return
}

func (s *server) GetCoach(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if isCoach, _ := s.checkCoachExistence(write, request, false); !isCoach {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}
	coachIdString := request.URL.Query().Get("coachId")
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	if p, ok_ := s.CoachesCache.ReadById(coachId); ok_ {
		log.Printf("coach loaded from cache: %d", (*p).Key)
		write.WriteHeader(http.StatusOK)
		WriteDataResponse(write, "Тренер получен", false, http.StatusOK, *p)
		return
	}
	coach, errGetCoach := s.Store.GetCoach(coachId)
	if errGetCoach != nil {
		log.Printf("error returns new coach data: %v\n", errGetCoach)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось получить тренера", true, http.StatusNotFound)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Тренер получен", false, http.StatusOK, coach)
	return
}

func (s *server) GetCoachFull(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if isCoach, _ := s.checkCoachExistence(write, request, false); !isCoach {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}
	coachIdString := request.URL.Query().Get("coachId")
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	coach, errGetCoach := s.Store.GetCoachFull(coachId)
	if errGetCoach != nil {
		log.Printf("error returns new coach data: %v\n", errGetCoach)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось получить тренера", true, http.StatusNotFound)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Тренер получен", false, http.StatusOK, coach)
	return
}

func (s *server) UpdateCoach(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if isCoach, _ := s.checkCoachExistence(write, request, false); !isCoach {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}
	decoder := json.NewDecoder(request.Body)
	var decoded requests.UpdateCoach
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
	coachIdString := request.URL.Query().Get("coachId")
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	sql := utils.GenerateUpdateSql("coach", coachId, params, values)
	errGetCoach := s.Store.UpdateCoach(sql)
	if errGetCoach != nil {
		log.Printf("error returns new coach data: %v\n", errGetCoach)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка обновления тренера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteResponse(write, "Тренер обновлен", false, http.StatusOK)
	return
}

func (s *server) DeleteCoach(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	coachIdString := request.URL.Query().Get("coachId")
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	errGetCoach := s.Store.DeleteCoach(coachId)
	if errGetCoach != nil {
		log.Printf("error returns new coach data: %v\n", errGetCoach)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось удалить тренера", true, http.StatusNotFound)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteResponse(write, "Тренер удален", false, http.StatusOK)
	return
}

func (s *server) GetAllPupilsForCoach(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if isCoach, _ := s.checkCoachExistence(write, request, false); !isCoach {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}
	coachIdString := request.URL.Query().Get("coachId")
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	pupils, errGetPupils := s.Store.GetCoachPupils(coachId)
	if errGetPupils != nil {
		log.Printf("error returns coach pupils list: %v\n", errGetPupils)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось получить список учениц", true, http.StatusNotFound)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список учениц получен", false, http.StatusOK, pupils)
	return
}

func (s *server) GetAllCoachList(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	coaches, errGetCoaches := s.Store.GetAllCoaches()
	if errGetCoaches != nil {
		log.Printf("cannot list admins: %v\n", errGetCoaches)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка сервера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список тренеров получен", false, http.StatusOK, coaches)
	return
}

func (s *server) GetNearestBirthdays(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if isCoach, _ := s.checkCoachExistence(write, request, false); !isCoach {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}
	coachIdString := request.URL.Query().Get("coachId")
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	birthDayList, errGetBdayList := s.Store.GetBirthdaysList(coachId)
	if errGetBdayList != nil {
		log.Printf("cannot list bdays for id: %d: %v\n", coachId, errGetBdayList)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка сервера", true, http.StatusInternalServerError)
		return
	}
	bdays := utils.GetNearestBdays(birthDayList)
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список ближайших дней рождений получен", false, http.StatusOK, bdays)
	return
}

func (s *server) ArchiveCoach(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	coachIdString := request.URL.Query().Get("coachId")
	if isAdmin, _ := s.checkExistence(write, request); !isAdmin {
		return
	}
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	if errArchiveCoach := s.Store.ArchiveCoach(coachId); errArchiveCoach != nil {
		log.Printf("error returns archive coach: %v\n", errArchiveCoach)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось архивировать тренера", true, http.StatusNotFound)
		return
	}
	s.CoachesCache.RemoveCoach(coachId)
	write.WriteHeader(http.StatusOK)
	WriteResponse(write, "Тренер архивирован", false, http.StatusOK)
	return
}
func (s *server) ArchiveCoachGet(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	coaches, errGetArchiveCoaches := s.Store.ArchiveCoachGet()
	if errGetArchiveCoaches != nil {
		log.Printf("error returns archive coach: %v\n", errGetArchiveCoaches)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Не удалось получить архивных тренеров", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, " Архивированные тренеры получены", false, http.StatusOK, coaches)
	return
}

func (s *server) DearchiveCoach(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	coachIdString := request.URL.Query().Get("coachId")
	if isAdmin, _ := s.checkExistence(write, request); !isAdmin {
		return
	}
	coachId, errConvCoach := strconv.Atoi(coachIdString)
	if errConvCoach != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	if errDearchivePupil := s.Store.DearchiveCoach(coachId); errDearchivePupil != nil {
		log.Printf("error returns archive pupil: %v\n", errDearchivePupil)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Не удалось разархивировать тренера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteResponse(write, "Тренер разархивирован", false, http.StatusOK)
	return
}
