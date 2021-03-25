package http

import (
	"DBproject/internal/service"
	"github.com/labstack/echo/v4"
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

// Безвозвратное удаление всей пользовательской информации из базы данных.
func (h *Handler) Clear(c echo.Context) error {
	return nil // todo
}

// Получение инфомарции о базе данных.
func (h *Handler) Status(c echo.Context) error {
	return nil // todo
}

