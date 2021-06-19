package repository

import (
	"DBproject/internal/service"
	"github.com/jackc/pgx"
)

type serviceRepo struct {
	DB *pgx.ConnPool
}

func NewServiceRepo(db *pgx.ConnPool) service.ServiceRepo {
	return &serviceRepo{
		DB: db,
	}
}

func (s serviceRepo) Clear() error {
	_, err := s.DB.Exec("truncate users, forums, threads, votes, posts, forum_users cascade")
	return err
}

func (s serviceRepo) Status() service.DBinfo {
	info := new(service.DBinfo)
	_ = s.DB.QueryRow("select count(*) from forums").Scan(&info.Forums)
	_ = s.DB.QueryRow("select count(*) from users").Scan(&info.Users)
	_ = s.DB.QueryRow("select count(*) from posts").Scan(&info.Posts)
	_ = s.DB.QueryRow("select count(*) from threads").Scan(&info.Threads)

	return *info
}
