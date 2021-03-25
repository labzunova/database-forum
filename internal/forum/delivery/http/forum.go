package http

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	ForumUcase forum.ForumUsecase
}

func NewForumHandler(forumUcase forum.ForumUsecase) forum.ForumHandler {
	handler := &Handler{
		ForumUcase: forumUcase,
	}

	return handler
}

type UsersParse struct {
	limit int32
	since string
	desc bool
}

// Создание нового форума.
func (h Handler) ForumCreate(c echo.Context) error {
	newForum := new(models.Forum)
	if err := c.Bind(newForum); err != nil {
		return err
	}

	forum, err := h.ForumUcase.CreateNewForum(*newForum)
	switch err.Message {
	case "404":
		return c.JSON(http.StatusNotFound, "Владелец форума не найден")
	case "409":
		return c.JSON(http.StatusConflict, forum)
	}

	return c.JSON(http.StatusCreated, forum)
}

// Получение информации о форуме по его идентификаторе
func (h Handler) ForumGetOne(c echo.Context) error {
	slug := c.Param("slug")

	forum, err := h.ForumUcase.GetForum(slug)
	if err.Message != "" {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, forum)
}

// Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
// Пользователи выводятся отсортированные по nickname в порядке возрастания.
// Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.
func (h Handler) ForumGetUsers(c echo.Context) error {
	slug := c.Param("slug")
	parametres := c.QueryParams() // TODO ?

	users, err := h.ForumUcase.GetUsers(slug, parametres)
	if err.Message != "" {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, users)
}

// Получение списка ветвей обсужления данного форума.
// Ветви обсуждения выводятся отсортированные по дате создания.
func (h Handler) ForumGetThreads(c echo.Context) error {
	slug := c.Param("slug")
	parametres := c.QueryParams() // TODO ?

	threads, err := h.ForumUcase.GetThreads(slug, parametres)
	if err.Message != "" {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, threads)
}

