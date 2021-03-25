package repository

import (
	"DBproject/internal/service"
	"database/sql"
)

type serviceRepo struct {
	DB *sql.DB
}

func NewServiceRepo(db *sql.DB) service.ServiceRepo {
	return &serviceRepo{
		DB: db,
	}
}

func (s serviceRepo) Clear() error {
	return nil // todo
}

func (s serviceRepo) Status() error {
	return nil // todo
}
