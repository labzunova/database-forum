package forum

import (
	"DBproject/internal/forum/delivery/http"
	"DBproject/models"
	"github.com/labstack/echo/v4"
)

type ForumHandler interface {
	ForumCreate(c echo.Context) error
 	ForumGetOne(c echo.Context) error
	ForumGetUsers(c echo.Context) error
	ForumGetThreads(c echo.Context) error
}

type ForumUsecase interface {
	CreateNewForum(forum models.Forum) (models.Forum, models.Error)
	GetForum(id string) (models.Forum, models.Error)
	GetUsers(slug string, params http.UsersParse) ([]models.User, models.Error)
	GetThreads(slug string, params http.UsersParse) ([]models.Thread, models.Error)
}

type ForumRepo interface {
	CreateNewForum(forum models.Forum) (models.Forum, models.Error)
	GetForum(id string) (models.Forum, models.Error)
	GetUsers(slug string, params http.UsersParse) ([]models.User, models.Error)
	GetThreads(slug string, params http.UsersParse) ([]models.Thread, models.Error)
}
