package threads

import (
	"DBproject/models"
	"github.com/labstack/echo/v4"
)

type ThreadsHandler interface {
	ThreadCreate(c echo.Context) error
	ThreadGetOne(c echo.Context) error
	ThreadUpdate(c echo.Context) error
	ThreadGetPosts(c echo.Context) error
	ThreadVote(c echo.Context) error
}

type ThreadsUsecase interface {
	CreateThread(slug string, thread models.Thread) (models.Thread, models.Error)
	GetThread(slug string) (models.Thread, models.Error)
	UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error)
	GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
	VoteThread(slug string, vote models.Vote) (models.Thread, models.Error)
}

type ThreadsRepo interface {
	CreateThread(slug string, thread models.Thread) (models.Thread, models.Error)
	GetThread(slug string) (models.Thread, models.Error)
	UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error)
	GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
	VoteThread(slug string, vote models.Vote) (models.Thread, models.Error)
}