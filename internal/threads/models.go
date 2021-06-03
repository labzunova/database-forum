package threads

import (
	"DBproject/models"
	"github.com/labstack/echo/v4"
)

type ThreadsHandler interface {
	ThreadGetOne(c echo.Context) error
	ThreadUpdate(c echo.Context) error
	ThreadGetPosts(c echo.Context) error
	ThreadVote(c echo.Context) error
}

type ThreadsUsecase interface {
	GetThread(slug string) (models.Thread, models.Error)
	UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error)
	GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
	VoteThread(slug string, vote models.Vote) (models.Thread, models.Error)
	SlugOrID(slug string) int
}

type ThreadsRepo interface {
	GetThread(slug string, id int) (models.Thread, models.Error)
	UpdateThreadBySlug(slug string, thread models.Thread) (models.Thread, models.Error)
	UpdateThreadById(id int, thread models.Thread) (models.Thread, models.Error)
	GetThreadPostsBySlug(slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
	GetThreadPostsById(id int, params models.ParseParamsThread) ([]models.Post, models.Error)
	VoteThreadBySlug(slug string, vote models.Vote) models.Error
	VoteThreadById(id int, vote models.Vote) models.Error
	UpdateVoteThreadBySlug(slug string, vote models.Vote) models.Error
	UpdateVoteThreadById(id int, vote models.Vote) models.Error
}
