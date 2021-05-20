package http

import (
	"DBproject/internal/posts"
	"DBproject/internal/service"
	"DBproject/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	PostsUcase posts.PostsUsecase
}

func NewPostsHandler(postsUcase posts.PostsUsecase) posts.PostsHandler {
	handler := &Handler{
		PostsUcase: postsUcase,
	}

	return handler
}

type relatedForPost struct {
}

// Получение информации о ветке обсуждения по его имени.
func (h *Handler) PostGetOne(c echo.Context) error {
	id := c.Param("id")
	related := c.QueryParam("related") // TODO ???

	posts, err :=  h.PostsUcase.GetPost(id, related) // todo
	if err.Message == "404" {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, posts)
}

// Изменение сообщения на форуме.
// Если сообщение поменяло текст, то оно должно получить отметку `isEdited`.
func (h *Handler) PostUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 0, 64)
	var newMessage posts.UpdateMessage
	if err := c.Bind(newMessage); err != nil {
		return err
	}

	post, err := h.PostsUcase.UpdatePost(id, newMessage.Message)
	if err.Message == "404" {
		return c.JSON(http.StatusNotFound, "Сообщение отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, post)
}

// Добавление новых постов в ветку обсуждения на форум.
// Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).
func (h *Handler) PostsCreate(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newPosts := make([]models.Post, 0)
	if err := c.Bind(newPosts); err != nil {
		return err
	}

	posts, err := h.PostsUcase.CreatePosts(slug, newPosts)
	switch err.Message {
	case "404":
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в базе")
	case "409":
		return c.JSON(http.StatusConflict, "Хотя бы один пост отсутствует в ветке")
	}

	return c.JSON(http.StatusOK, posts)
}


