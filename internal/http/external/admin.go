package external

import (
	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"dnevnik-rg.ru/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *server) CreateAdmin(write http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		write.WriteHeader(http.StatusNotFound)
		WriteResponse(write, "Неизвестный метод", true, http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(request.Body)
	var decoded requests.NewAdmin
	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		write.WriteHeader(http.StatusBadRequest)
		WriteResponse(write, "Ошибка валидации", true, http.StatusBadRequest)
		return
	}
	key := int(utils.GetUdid())
	newAdmin := models.Admin{
		General: models.General{Key: key,
			Fio:     decoded.Fio,
			DateReg: time.Now().Format(time.RFC3339),
			LogoUri: "https://dnevnik-rg.ru/admin-logo.png",
			Role:    "ADMIN"},
	}
	newPassword := models.Password{
		Key:        key,
		CheckSum:   utils.NewPassword(),
		LastUpdate: time.Now().Format(time.RFC3339),
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
			write.WriteHeader(http.StatusInternalServerError)
			WriteResponse(write, "Ошибка создания администратора", true, http.StatusInternalServerError)
			return
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

func (s *server) GetAdmin() {

}
