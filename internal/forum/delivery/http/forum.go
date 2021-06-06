package http

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"fmt"
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
		errr.Message = "Can't find user with nickname: " + newForum.User
		return c.JSON(http.StatusNotFound, errr)
	case 409:
		forumOld, _ := h.ForumUcase.GetForum(newForum.Slug)
		return c.JSON(http.StatusConflict, forumOld)
	}

	return c.JSON(http.StatusCreated, forum)
}

// ForumGetOne Получение информации о форуме по его идентификаторе
func (h Handler) ForumGetOne(c echo.Context) error {
	slug := c.Param("slug")

	forumResponse, err := h.ForumUcase.GetForum(slug)
	if err.Code != 200 {
		//return c.JSON(http.StatusOK, nil) // TODO FOR PERF TESTS ????
		//err.Message = "Can't find forum with slug: " + slug
		err.Message = "Can't find forum"
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, forumResponse)
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
		return c.JSON(http.StatusNotFound, err)
	case 409:
		return c.JSON(http.StatusConflict, thread)
	}

	return c.JSON(http.StatusCreated, thread)
}

// ForumGetUsers Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
// Пользователи выводятся отсортированные по nickname в порядке возрастания.
// Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.
func (h Handler) ForumGetUsers(c echo.Context) error {
	slug := c.Param("slug")
	getUsers := new(models.ParseParams)

	getUsers.Limit, _ = strconv.Atoi(c.QueryParam("limit"))
	getUsers.Since = c.QueryParam("since")
	getUsers.Desc, _ = strconv.ParseBool(c.QueryParam("desc"))

	users, err := h.ForumUcase.GetUsers(slug, *getUsers)
	if err.Code == 404 {
		//return c.JSON(http.StatusOK, nil) // TODO FOR PERF TESTS ????
		err.Message = "Can't find forum by slug: " + slug
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, users)
}

// ForumGetThreads Получение списка ветвей обсужления данного форума.
// Ветви обсуждения выводятся отсортированные по дате создания.
func (h Handler) ForumGetThreads(c echo.Context) error {
	slug := c.Param("slug")
	getThreads := new(models.ParseParams)

	getThreads.Limit, _ = strconv.Atoi(c.QueryParam("limit"))
	getThreads.Since = c.QueryParam("since")
	getThreads.Desc, _ = strconv.ParseBool(c.QueryParam("desc"))
	fmt.Println(getThreads)

	threads, err := h.ForumUcase.GetThreads(slug, *getThreads)
	if err.Code != 200 {
		//return c.JSON(http.StatusOK, nil) // TODO FOR PERF TESTS ????
		err.Message = "Can't find forum by slug: " + slug
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, threads)
}
