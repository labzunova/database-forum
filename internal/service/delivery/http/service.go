package http

import (
	"DBproject/helpers"
	"DBproject/internal/service"
	"net/http"
)

type Handler struct {
	serviceRepository service.ServiceRepo
}

func NewServiceHandler(repo service.ServiceRepo) service.ServiceHandler {
	handler := &Handler{
		serviceRepository: repo,
	}
	return handler
}

// Clear Безвозвратное удаление всей пользовательской информации из базы данных.
func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	err := h.serviceRepository.Clear()
	helpers.CreateResponse(w, http.StatusOK, err)
	return
}

// Status Получение инфомарции о базе данных.
func (h *Handler) Status(w http.ResponseWriter, r *http.Request)  {
	status := h.serviceRepository.Status()
	helpers.CreateResponse(w, http.StatusOK, status)
	return
}
