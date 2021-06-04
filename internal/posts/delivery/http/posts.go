package http

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
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

// PostGetOne Получение информации о ветке обсуждения по его имени.
	func (h *Handler) PostGetOne(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 0, 64)
	related := c.QueryParam("related")
	related = strings.ReplaceAll(related, "[", "")
	related = strings.ReplaceAll(related, "]", "")
	relatedSlice := strings.Split(related, ",")


	post, err := h.PostsUcase.GetPost(int(id))
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	fullPost, errr := h.PostsUcase.GetPostInfo(int(id), relatedSlice)
	if errr.Code != 200 {
		return c.JSON(http.StatusInternalServerError, err.Code)
	}

	fullPost.Post = post

	return c.JSON(http.StatusOK, fullPost)
}

// PostUpdate Изменение сообщения на форуме.
// Если сообщение поменяло текст, то оно должно получить отметку `isEdited`.
func (h *Handler) PostUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 0, 64)
	var newMessage posts.UpdateMessage
	if err := c.Bind(newMessage); err != nil {
		return err
	}

	post, err := h.PostsUcase.UpdatePost(int(id), newMessage.Message)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Сообщение отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, post)
}

// PostsCreate Добавление новых постов в ветку обсуждения на форум.
// Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).
func (h *Handler) PostsCreate(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newPosts := make([]models.Post, 0)
	//newPosts := models.PostsToCreate{}
	if err := c.Bind(&newPosts); err != nil {
		return c.JSON(http.StatusCreated, newPosts)
	}
	fmt.Println("AAAAAAAAAAAAAA")

	posts, err := h.PostsUcase.CreatePosts(slug, newPosts)
	switch err.Code {
	case 404:
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в базе")
	case 409:
		return c.JSON(http.StatusConflict, "Хотя бы один пост отсутствует в ветке")
	}

	return c.JSON(http.StatusOK, posts)
}
