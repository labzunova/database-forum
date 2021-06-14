package service

import (
	"net/http"
)

type DBinfo struct {
	Users   int `json:"user"`
	Forums  int `json:"forum"`
	Threads int `json:"thread"`
	Posts   int `json:"post"`
}

type ServiceHandler interface {
	Clear(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
}

type ServiceUsecase interface {
	Clear() error
	Status() DBinfo
}

type ServiceRepo interface {
	Clear() error
	Status() DBinfo
}
