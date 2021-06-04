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

// ForumCreate Создание нового форума.
func (h Handler) ForumCreate(c echo.Context) error {
	newForum := new(models.Forum)
	err := c.Bind(newForum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "")
	}

	forum, errr := h.ForumUcase.CreateNewForum(*newForum)
	switch errr.Code {
	case 404:
		return c.JSON(http.StatusNotFound, "Владелец форума не найден")
	case 409:
		forumOld, _ := h.ForumUcase.GetForum(forum.Slug)
		return c.JSON(http.StatusConflict, forumOld)
	}

	return c.JSON(http.StatusCreated, forum)
}

// ForumGetOne Получение информации о форуме по его идентификаторе
func (h Handler) ForumGetOne(c echo.Context) error {
	slug := c.Param("slug")

	forum, err := h.ForumUcase.GetForum(slug)
	if err.Code != 200 {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, forum)
}

// ThreadCreate Добавление новой ветки обсуждения на форум
func (h Handler) ThreadCreate(c echo.Context) error {
	slug := c.Param("slug")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		return err
	}

	thread, err := h.ForumUcase.CreateThread(slug, *newThread)
	switch err.Code {
	case 404:
		return c.JSON(http.StatusNotFound, "Автор ветки или форум не найдены")
	case 409:
		return c.JSON(http.StatusConflict, thread)
	}

	return c.JSON(http.StatusOK, thread)
}

// ForumGetUsers Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
// Пользователи выводятся отсортированные по nickname в порядке возрастания.
// Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.
func (h Handler) ForumGetUsers(c echo.Context) error {
	slug := c.Param("slug")
	getUsers := new(models.ParseParams)

	var err error
	getUsers.Limit, err = strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "limit не найден")
	}

	getUsers.Since = c.QueryParam("since")

	getUsers.Desc, err = strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "desc не найден")
	}

	users, errr := h.ForumUcase.GetUsers(slug, *getUsers)
	if errr.Code == 400 {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, users)
}

// ForumGetThreads Получение списка ветвей обсужления данного форума.
// Ветви обсуждения выводятся отсортированные по дате создания.
func (h Handler) ForumGetThreads(c echo.Context) error {
	slug := c.Param("slug")
	getThreads := new(models.ParseParams)

	var err error
	getThreads.Limit, err = strconv.Atoi(c.QueryParam("limit")) // todo if zero?
	if err != nil {
		return c.JSON(http.StatusBadRequest, "limit не найден")
	}

	getThreads.Since = c.QueryParam("since")

	getThreads.Desc, err = strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "desc не найден")
	}

	threads, errr := h.ForumUcase.GetThreads(slug, *getThreads)
	if errr.Code != 200 {
		return c.JSON(http.StatusNotFound, "Форум отсутствует в системе")
	}

	return c.JSON(http.StatusOK, threads)
}

