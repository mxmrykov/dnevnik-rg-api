package external

import (
	"net/http"

	"dnevnik-rg.ru/internal/models"
)

func (s *server) ShowCacheUsers(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	ok, _ := s.checkExistence(write, request)
	if !ok {
		return
	}
	cache := struct {
		Pupils  map[int]*models.Pupil `json:"pupils"`
		Coaches map[int]*models.Coach `json:"coaches"`
		Admins  map[int]*models.Admin `json:"admins"`
	}{
		s.PupilsCache.ReadAll(),
		s.CoachesCache.ReadAll(),
		s.AdminsCache.ReadAll(),
	}
	write.WriteHeader(http.StatusOK)
	WriteDataResponse(write, "Тренер получен", false, http.StatusOK, cache)
	return
}
