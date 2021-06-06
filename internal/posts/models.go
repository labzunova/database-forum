package posts

import (
	"DBproject/models"
	"github.com/labstack/echo/v4"
)

type UpdateMessage struct {
	Message string `json:"message"`
}

type FullPost struct {
	Post   *models.Post   `json:"post"`
	User   *models.User   `json:"author"`
	Forum  *models.Forum  `json:"forum"`
	Thread *models.Thread `json:"thread"`
}

type PostsHandler interface {
	PostGetOne(c echo.Context) error
	PostUpdate(c echo.Context) error
	PostsCreate(c echo.Context) error
}

type PostsUsecase interface {
	GetPost(id int) (models.Post, models.Error)
	GetPostInfo(id int, related []string) (FullPost, models.Error)

	UpdatePost(id int, message string) (models.Post, models.Error)
	CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error)
}

type PostsRepo interface {
	GetPost(id int) (models.Post, models.Error)
	GetPostAuthor(pid int) (models.User, models.Error)
	GetPostThread(pid int) (models.Thread, models.Error)
	GetPostForum(pid int) (models.Forum, models.Error)

	UpdatePost(id int, message string) (models.Post, models.Error)
	CreatePosts(thread models.Thread, posts []models.Post) ([]models.Post, models.Error)
	GetThreadAndForumById(id int) (models.Thread, models.Error)
	GetThreadAndForumBySlug(slug string) (models.Thread, models.Error)
	CheckValidParent(id, parent int) bool
}
