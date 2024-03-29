package external

import (
	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"dnevnik-rg.ru/pkg/utils"
)

func (s *server) GetCoachSchedule(write http.ResponseWriter, request *http.Request) {
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
	classDateString := request.URL.Query().Get("date")

	if _, errParse := time.Parse("2006-01-02", classDateString); errParse != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный формат времени", true, http.StatusBadRequest)
		return
	}
	schedule, errGetSchedule := s.Repository.GetCoachSchedule(coachId, classDateString)
	if errGetSchedule != nil {
		log.Printf("error returns coach pupils list: %v\n", errGetSchedule)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Не удалось получить расписание на тренера", true, http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Расписание на тренера получено", false, http.StatusOK, utils.GetAvailClassesTimesAlgo(schedule))
	return
}

func (s *server) CreateClass(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
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
	var decoded requests.CreateClass
	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}

	var (
		err        error
		class      requests.CreateClass
		newClassId int
	)

	if decoded.ClassType != "SINGLE" && decoded.ClassType != "MULTIPLE" {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неизвестный тип занятия", true, http.StatusBadRequest)
		return
	}

	if decoded.ClassTime == "SINGLE" && len(decoded.Pupil) != 1 {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "На одиночном занятии может присутствовать только одна ученица", true, http.StatusBadRequest)
		return
	}

	class.ClassType = decoded.ClassType
	class.Pupil = decoded.Pupil
	class.Coach = decoded.Coach

	if _, err = time.Parse("2006-01-02", decoded.ClassDate); err != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный формат времени", true, http.StatusBadRequest)
		return
	}

	class.ClassDate = decoded.ClassDate

	if _, err = time.Parse("15:04", decoded.ClassTime); err != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный формат времени", true, http.StatusBadRequest)
		return
	}

	class.ClassTime = decoded.ClassTime

	if _, err = time.Parse("15:04", decoded.Duration); err != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный формат длительности", true, http.StatusBadRequest)
		return
	}

	class.Duration = decoded.Duration

	if decoded.Capacity < 1 {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверное количество допустимых учениц", true, http.StatusBadRequest)
		return
	}

	class.Capacity = decoded.Capacity

	if decoded.Price < 1 {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Слишком маленькая цена", true, http.StatusBadRequest)
		return
	}

	class.Price = decoded.Price

	newClassId, err = s.Repository.CreateClass(class)
	if err != nil {
		log.Println("err at creating class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	if newClassId == -1 {
		write.WriteHeader(http.StatusConflict)
		WriteResponse(write, "Занятие на это время уже занято", true, http.StatusConflict)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Занятие успешно создан", false, http.StatusOK, newClassId)
	return
}

func (s *server) GetClassesTodayAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}

	if ok, _ := s.checkExistence(write, request); !ok {
		return
	}

	var (
		err  error
		date string
		res  []models.ShortStringClassInfo
	)

	date = request.URL.Query().Get("date")

	classes, err := s.Repository.GetAdminClassesForToday(date)

	if err != nil {
		log.Println("err at creating class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	for _, class := range classes {
		var pupils []int
		for _, pupilID := range class.Pupils {
			pupils = append(pupils, pupilID)
		}
		stringPupils, err := s.Repository.GetPupilsNameByIds(pupils)
		if err != nil {
			log.Println("err at getting pupils names:", err)
			write.WriteHeader(http.StatusInternalServerError)
			WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
			return
		}

		response := models.ShortStringClassInfo{
			Key:           class.Key,
			Coach:         class.Coach,
			Pupils:        make([]string, 0, len(pupils)),
			ClassTime:     class.ClassTime,
			ClassDuration: class.ClassDuration,
		}

		for _, pupil := range stringPupils {
			response.Pupils = append(response.Pupils, pupil)
		}

		res = append(res, response)
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список занятий на сегодня получен", false, http.StatusOK, res)
	return
}
