package external

import (
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
	classDateTime, errParse := time.Parse("2006-01-02", classDateString)
	if errParse != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Неверный формат времени", true, http.StatusBadRequest)
		return
	}
	classDateFormatedString := classDateTime.Format(time.RFC3339)
	schedule, errGetSchedule := s.Repository.GetCoachSchedule(coachId, classDateFormatedString)
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
