package posts

import (
	"DBproject/models"
	"github.com/labstack/echo/v4"
)

type PostsHandler interface {
	PostGetOne(c echo.Context) error
	PostUpdate(c echo.Context) error
	PostsCreate(c echo.Context) error
}

type PostsUsecase interface {
	GetPost() ([]models.Post, models.Error)
	UpdatePost(id int64, post models.Post) (models.Post, models.Error)
	CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error)
}


type PostsRepo interface {
	GetPost() ([]models.Post, models.Error)
	UpdatePost(id int64, post models.Post) (models.Post, models.Error)
	CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error)
}
