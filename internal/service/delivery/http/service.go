package http

import (
	"DBproject/helpers"
	"DBproject/internal/service"
	"net/http"
)

type Handler struct {
	ServiceUcase service.ServiceUsecase
}

func NewServiceHandler(serviceUcase service.ServiceUsecase) service.ServiceHandler {
	handler := &Handler{
		ServiceUcase: serviceUcase,
	}
	return handler
}

// Clear Безвозвратное удаление всей пользовательской информации из базы данных.
func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	err := h.ServiceUcase.Clear()
	helpers.CreateResponse(w, http.StatusOK, err)
	return
}

// Status Получение инфомарции о базе данных.
func (h *Handler) Status(w http.ResponseWriter, r *http.Request)  {
	status := h.ServiceUcase.Status()
	helpers.CreateResponse(w, http.StatusOK, status)
	return
}
