package http

import (
	"DBproject/internal/service"
	"github.com/labstack/echo/v4"
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
func (h *Handler) Clear(c echo.Context) error {
	err := h.ServiceUcase.Clear()
	return c.JSON(http.StatusOK, err)
}

// Status Получение инфомарции о базе данных.
func (h *Handler) Status(c echo.Context) error {
	status := h.ServiceUcase.Status()
	return c.JSON(http.StatusOK, status)
}
