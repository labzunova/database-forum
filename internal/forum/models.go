package forum

import (
	"DBproject/models"
	"net/http"
)

type ForumHandler interface {
	ForumCreate(w http.ResponseWriter, r *http.Request)
	ForumGetOne(w http.ResponseWriter, r *http.Request)
	ThreadCreate(w http.ResponseWriter, r *http.Request)
	ForumGetUsers(w http.ResponseWriter, r *http.Request)
	ForumGetThreads(w http.ResponseWriter, r *http.Request)
}

//type ForumUsecase interface {
//	CreateNewForum(forum models.Forum) (models.Forum, models.Error)
//	GetForum(id string) (models.Forum, models.Error)
//	CreateThread(slug string, thread models.Thread) (models.Thread, models.Error)
//	GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error)
//	GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error)
//}

type ForumRepo interface {
	CreateNewForum(forum models.Forum) (models.Forum, models.Error)
	GetForum(id string) (models.Forum, models.Error)
	CreateThread(slug string, thread models.Thread) (models.Thread, models.Error)
	GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error)
	GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error)

	GetThreadBySlug(slug string) (models.Thread, models.Error)
	CheckForumExists(slug string) bool
}
