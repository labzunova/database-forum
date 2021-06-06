package http

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	ThreadsUcase threads.ThreadsUsecase
}

func NewThreadsHandler(threadsUcase threads.ThreadsUsecase) threads.ThreadsHandler {
	handler := &Handler{
		ThreadsUcase: threadsUcase,
	}
	return handler
}

// ThreadGetOne Получение информации о ветке обсуждения по его имени.
func (h *Handler) ThreadGetOne(c echo.Context) error {
	slug := c.Param("slug_or_id")

	thread, err :=  h.ThreadsUcase.GetThread(slug)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, err)
	}

	//thread.Created = thread.Created.Add(-time.Hour * 3) // TODO ВРЕМЕННО ДЛЯ КОМПА

	return c.JSON(http.StatusOK, thread)
}

// ThreadUpdate Обновление ветки обсуждения на форуме.
func (h *Handler) ThreadUpdate(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		return err
	}

	thread, err :=  h.ThreadsUcase.UpdateThread(slugOrId, *newThread)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, err)
	}

	//thread.Created = thread.Created.Add(-time.Hour * 3) // TODO ВРЕМЕННО ДЛЯ КОМПА

	return c.JSON(http.StatusOK, thread)
}

// ThreadGetPosts Получение списка сообщений в данной ветке форуме.
// Сообщения выводятся отсортированные по дате создания.
func (h *Handler) ThreadGetPosts(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	getPosts := new(models.ParseParamsThread)

	getPosts.Limit, _ = strconv.Atoi(c.QueryParam("limit"))
	getPosts.Since, _ = strconv.Atoi(c.QueryParam("since"))
	getPosts.Sort = c.QueryParam("sort")
	getPosts.Desc, _ = strconv.ParseBool(c.QueryParam("desc"))

	posts, err := h.ThreadsUcase.GetThreadPosts(slugOrId, *getPosts)
	if err.Code == 404 {
		err.Message = "Can't find forum by slug: " + slugOrId
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, posts)
}

// ThreadVote Изменение голоса за ветвь обсуждения.
// Один пользователь учитывается только один раз и может изменить своё мнение.
func (h *Handler) ThreadVote(c echo.Context) error {
	fmt.Println("VOTE handler")
	newVote := new(models.Vote)
	if err := c.Bind(newVote); err != nil {
		return err
	}
	slug := c.Param("slug_or_id")

	thread, err :=  h.ThreadsUcase.VoteThread(slug, *newVote)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, err)
	}
	if err.Code != 200 {
		return c.JSON(http.StatusInternalServerError, "Error")
	}

	//thread.Created = thread.Created.Add(-time.Hour * 3) // TODO ВРЕМЕННО ДЛЯ КОМПА

	return c.JSON(http.StatusOK, thread)
}
