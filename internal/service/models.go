package service

import (
	"github.com/labstack/echo/v4"
)

type ServiceHandler interface {
	Clear(c echo.Context) error
	Status(c echo.Context) error
}

type ServiceUsecase interface {
	Clear() error
	Status() error
}

type ServiceRepo interface {
	Clear() error
	Status() error
}

