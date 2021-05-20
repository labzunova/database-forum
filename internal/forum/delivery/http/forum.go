package http

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

// Создание нового форума.
func (h Handler) ForumCreate(c echo.Context) error {
	newForum := new(models.Forum)
	err := c.Bind(newForum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	forum, errr := h.ForumUcase.CreateNewForum(*newForum)
	switch errr.Message {
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
	getUsers := new(models.ParseParams)

	limit, err := strconv.Atoi(c.QueryParam("limit")) // todo if zero?
	if err != nil {
		// todo
	}
	getUsers.Limit = int32(limit)

	var since string
	since = c.QueryParam("since") // todo if zero?
	getUsers.Since = since

	var desc bool
	desc, err = strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		// todo
	}
	getUsers.Desc = desc

	users, errr := h.ForumUcase.GetUsers(slug, *getUsers)
	if errr.Message != "" {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, users)
}

// Получение списка ветвей обсужления данного форума.
// Ветви обсуждения выводятся отсортированные по дате создания.
func (h Handler) ForumGetThreads(c echo.Context) error {
	slug := c.Param("slug")
	getThreads := new(models.ParseParams)

	limit, err := strconv.Atoi(c.QueryParam("limit")) // todo if zero?
	if err != nil {
		// todo
	}
	getThreads.Limit = int32(limit)

	var since string
	since = c.QueryParam("since") // todo if zero?
	getThreads.Since = since

	var desc bool
	desc, err = strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		// todo
	}
	getThreads.Desc = desc

	threads, errr := h.ForumUcase.GetThreads(slug, *getThreads)
	if errr.Message != "" {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, threads)
}

