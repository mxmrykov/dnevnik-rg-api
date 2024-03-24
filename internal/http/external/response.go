package external

import (
	"dnevnik-rg.ru/internal/models/response"
	"encoding/json"
	"net/http"
)

func WriteResponse(write http.ResponseWriter, message string, isError bool, statusCode int) {
	resp, errMarshaling := json.MarshalIndent(response.Response{
		Data:       nil,
		StatusCode: statusCode,
		Message:    message,
		IsError:    isError,
	}, "", "\t")
	if errMarshaling != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = write.Write(resp)
}

func WriteDataResponse(write http.ResponseWriter, message string, isError bool, statusCode int, data interface{}) {
	resp, errMarshaling := json.MarshalIndent(response.Response{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
		IsError:    isError,
	}, "", "\t")
	if errMarshaling != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = write.Write(resp)
}
