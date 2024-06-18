package external

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/internal/models/response"
	"github.com/golang-jwt/jwt"

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
	schedule, errGetSchedule := s.Store.GetCoachSchedule(coachId, classDateString)
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

	newClassId, err = s.Store.CreateClass(class)
	if err != nil {
		log.Println("err at creating class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	if newClassId == -1 {
		log.Println("err at creating class: class at this time exists exists")
		write.WriteHeader(http.StatusConflict)
		WriteResponse(write, "Занятие на это время уже занято", true, http.StatusConflict)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Занятие успешно создано", false, http.StatusOK, newClassId)
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

	classes, err := s.Store.GetAdminClassesForToday(date)

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
		stringPupils, err := s.Store.GetPupilsNameByIds(pupils)
		if err != nil {
			log.Println("err at getting pupils names:", err)
			write.WriteHeader(http.StatusInternalServerError)
			WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
			return
		}
		coachString, err := s.Store.GetCoachNameById(class.Coach)
		if err != nil {
			log.Println("err at getting pupils names:", err)
			write.WriteHeader(http.StatusInternalServerError)
			WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
			return
		}

		resp := models.ShortStringClassInfo{
			Key:            class.Key,
			Coach:          *coachString,
			Pupils:         make([]string, 0, len(pupils)),
			ClassTime:      class.ClassTime,
			ClassDuration:  class.ClassDuration,
			ClassType:      class.ClassType,
			PupilCount:     class.PupilCount,
			Scheduled:      class.Scheduled,
			IsOpenToSignUp: class.IsOpenToSignUp,
		}

		for _, pupil := range stringPupils {
			resp.Pupils = append(resp.Pupils, pupil)
		}

		res = append(res, resp)
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список занятий на сегодня получен", false, http.StatusOK, res)
	return
}

func (s *server) GetClassesTodayCoach(write http.ResponseWriter, request *http.Request) {

}
func (s *server) GetClassesTodayPupil(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}

	Auth := request.Header.Get("Authorization")
	user, errConv := strconv.Atoi(Auth)

	if errConv != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный ID пользователя", true, http.StatusBadRequest)
		return
	}

	var date string

	date = request.URL.Query().Get("date")

	if !verifyPupil(request) {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}

	classes, err := s.Store.GetTodayPupilClasses(user, date)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Ошибка сервера", true, http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список занятий на сегодня получен", false, http.StatusOK, classes)
	return
}

func (s *server) CancelClass(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}

	classIdString := request.URL.Query().Get("classId")
	classId, errConv := strconv.Atoi(classIdString)

	if errConv != nil {
		log.Println("err converting classId from query:", errConv)
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusBadRequest)
		return
	}

	if err := s.Store.CancelClass(classId); err != nil {
		log.Println("err at canceling class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Занятие отменено", false, http.StatusOK, response.CanceledClass{Canceled: true, Key: classId})
	return
}

func (s *server) DeleteClass(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}

	if isCoach, _ := s.checkCoachExistence(write, request, false); !isCoach {
		if ok, _ := s.checkExistence(write, request); !ok {
			return
		}
	}

	classIdString := request.URL.Query().Get("classId")
	classId, errConv := strconv.Atoi(classIdString)

	if errConv != nil {
		log.Println("err converting classId from query:", errConv)
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusBadRequest)
		return
	}

	if err := s.Store.DeleteClass(classId); err != nil {
		log.Println("err at creating class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Занятие удалено", false, http.StatusOK, response.DeletedClass{Deleted: true, Key: classId})
	return
}

func (s *server) GetClassesMonthAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if ok, _ := s.checkExistence(write, request); !ok {
		return
	}
	today := time.Now()
	lastDate := time.Now().Add(24 * 30 * time.Hour)

	var month string

	if today.Month() < 10 {
		month = fmt.Sprintf("0%d", today.Month())
	} else {
		month = fmt.Sprintf("%d", today.Month())
	}

	stringToday := fmt.Sprintf("%d-%s-%d", today.Year(), month, today.Day())

	if lastDate.Month() < 10 {
		month = fmt.Sprintf("0%d", lastDate.Month())
	} else {
		month = fmt.Sprintf("%d", lastDate.Month())
	}

	stringLastDay := fmt.Sprintf("%d-%s-%d", lastDate.Year(), month, lastDate.Day())

	classes, err := s.Store.GetClassesForMonth("ADMIN", stringToday, stringLastDay)
	if err != nil {
		log.Println("err at creating class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	c, err := utils.FillMonthCalendar(stringToday, stringLastDay, classes)
	if err != nil {
		log.Println("err at listing classes:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список занятий получен", false, http.StatusOK, c)
	return

}
func (s *server) GetClassesMonthCoach(write http.ResponseWriter, request *http.Request) {
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
}
func (s *server) GetClassesMonthPupil(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
}

func (s *server) GetClassInfoAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	if ok, _ := s.checkExistence(write, request); !ok {
		return
	}

	classIdString := request.URL.Query().Get("classId")
	classId, err := strconv.Atoi(classIdString)
	if err != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный ID занятия", true, http.StatusBadRequest)
		return
	}

	class, err := s.Store.GetClassById(classId)

	if err != nil && err.Error() == "no rows in result set" {
		log.Println("cannot resolve class with such id:", classId)
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Занятия не существует", true, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println("err at getting class:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	var pupils []int
	for _, pupilID := range class.Pupils {
		pupils = append(pupils, pupilID)
	}
	stringPupils, err := s.Store.GetPupilsNameByIds(pupils)
	if err != nil {
		log.Println("err at getting pupils names:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}
	coachString, err := s.Store.GetCoachNameById(class.Coach)
	if err != nil {
		log.Println("err at getting pupils names:", err)
		write.WriteHeader(http.StatusInternalServerError)
		WriteResponse(write, "Произошла ошибка на сервере", true, http.StatusInternalServerError)
		return
	}

	resp := models.ShortStringClassInfo{
		Key:            class.Key,
		Coach:          *coachString,
		Pupils:         make([]string, 0, len(pupils)),
		ClassDate:      class.ClassDate,
		ClassTime:      class.ClassTime,
		ClassDuration:  class.ClassDuration,
		ClassType:      class.ClassType,
		PupilCount:     class.PupilCount,
		Scheduled:      class.Scheduled,
		IsOpenToSignUp: class.IsOpenToSignUp,
	}

	for _, pupil := range stringPupils {
		resp.Pupils = append(resp.Pupils, pupil)
	}

	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Список занятий получен", false, http.StatusOK, resp)
	return
}

func verifyPupil(request *http.Request) bool {
	Auth := request.Header.Get("Authorization")
	_, errParse := jwt.Parse(Auth, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if errParse != nil {
		log.Println("error at parsing token: ", errParse)
		return false
	}
	payload, err := utils.GetJwtPayload(Auth)
	if err != nil {
		log.Println("error at parsing token: ", errParse)
		return false
	}
	if convId, errConvId := strconv.Atoi(
		request.Header.Get("X-User-Id"),
	); errConvId == nil && convId == payload.Key {
		return true
	} else {
		return false
	}
}
