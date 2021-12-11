package threads

import (
	"DBproject/models"
	"net/http"
)

type ThreadsHandler interface {
	ThreadGetOne(w http.ResponseWriter, r *http.Request)
	ThreadUpdate(w http.ResponseWriter, r *http.Request)
	ThreadGetPosts(w http.ResponseWriter, r *http.Request)
	ThreadVote(w http.ResponseWriter, r *http.Request)
}

//type ThreadsUsecase interface {
//	GetThread(slug string) (models.Thread, models.Error)
//	UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error)
//	GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
//	VoteThread(slug string, vote models.Vote) (models.Thread, models.Error)
//	SlugOrID(slug string) int
//}

type ThreadsRepo interface {
	GetThread(slug string, id int) (models.Thread, models.Error)
	UpdateThreadBySlug(slug string, thread models.Thread) (models.Thread, models.Error)
	UpdateThreadById(id int, thread models.Thread) (models.Thread, models.Error)
	GetThreadPostsBySlug(slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
	GetThreadPostsById(id int, slug string, params models.ParseParamsThread) ([]models.Post, models.Error)
	VoteThreadBySlug(slug string, vote models.Vote) models.Error
	VoteThreadById(id int, vote models.Vote) models.Error

	GetThreadIDBySlug(slug string, id int) (int, models.Error)
}
