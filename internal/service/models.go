package service

import (
	"github.com/labstack/echo/v4"
)

type DBinfo struct {
	Users int `json:"user"`
	Forums int `json:"forum"`
	Threads int `json:"thread"`
	Posts int `json:"post"`
}

type ServiceHandler interface {
	Clear(c echo.Context) error
	Status(c echo.Context) error
}

type ServiceUsecase interface {
	Clear() error
	Status() DBinfo
}

type ServiceRepo interface {
	Clear() error
	Status() DBinfo
}

