package usecase

import (
	"DBproject/internal/service"
)

type serviceUsecase struct {
	serviceRepository service.ServiceRepo
}

func NewServiceUsecase(repo service.ServiceRepo) service.ServiceUsecase {
	return &serviceUsecase{
		serviceRepository: repo,
	}
}
func (s serviceUsecase) Clear() error {
	return s.serviceRepository.Clear()
}

func (s serviceUsecase) Status() service.DBinfo {
	return s.serviceRepository.Status()
}

