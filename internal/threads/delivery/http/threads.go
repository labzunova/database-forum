package http

import (
	"DBproject/internal/threads"
	"DBproject/models"
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

type ThreadsParse struct {
	limit int32
	since string
	sort string
	desc bool
}

func (h *Handler) ThreadCreate(c echo.Context) error {
	slug := c.Param("slug")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		return err
	}

	thread, err := h.ThreadsUcase.CreateThread(slug, *newThread)
	switch err.Message {
	case "404":
		return c.JSON(http.StatusNotFound, "Автор ветки или форум не найдены")
	case "409":
		return c.JSON(http.StatusConflict, thread)
	}

	return c.JSON(http.StatusOK, thread)
}

// Получение информации о ветке обсуждения по его имени.
func (h *Handler) ThreadGetOne(c echo.Context) error {
	slug := c.Param("slug_or_id")

	thread, err :=  h.ThreadsUcase.GetThread(slug)
	if err.Message == "404" {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)
}

// Обновление ветки обсуждения на форуме.
func (h *Handler) ThreadUpdate(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		return err
	}

	thread, err :=  h.ThreadsUcase.UpdateThread(slug, *newThread)
	if err.Message == "404" {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)
}

// Получение списка сообщений в данной ветке форуме.
// Сообщения выводятся отсортированные по дате создания.
func (h *Handler) ThreadGetPosts(c echo.Context) error {
	slug := c.Param("slug_or_id")
	limit, _ := strconv.ParseInt(c.QueryParam("limit"), 0, 32)
	since := c.QueryParam("since")
	sort := c.QueryParam("sort")
	desc, _ := strconv.ParseBool(c.QueryParam("desc"))
	params := ThreadsParse{limit, since, sort, desc }

	posts, err := h.ThreadsUcase.GetThreadPosts(slug, params)
	if err.Message == "404" {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, posts)
}

// Изменение голоса за ветвь обсуждения.
// Один пользователь учитывается только один раз и может изменить своё мнение.
func (h *Handler) ThreadVote(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newVote := new(models.Vote)
	if err := c.Bind(newVote); err != nil {
		return err
	}

	thread, err :=  h.ThreadsUcase.VoteThread(slug, *newVote)
	if err.Message == "404" {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)

}
