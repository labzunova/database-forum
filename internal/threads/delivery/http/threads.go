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

// ThreadGetOne Получение информации о ветке обсуждения по его имени.
func (h *Handler) ThreadGetOne(c echo.Context) error {
	slug := c.Param("slug_or_id")

	thread, err :=  h.ThreadsUcase.GetThread(slug)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)
}

// ThreadUpdate Обновление ветки обсуждения на форуме.
func (h *Handler) ThreadUpdate(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		return err
	}

	thread, err :=  h.ThreadsUcase.UpdateThread(slug, *newThread)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)
}

// Получение списка сообщений в данной ветке форуме.
// Сообщения выводятся отсортированные по дате создания.
func (h *Handler) ThreadGetPosts(c echo.Context) error {
	slug := c.Param("slug_or_id")
	getPosts := new(models.ParseParamsThread)

	limit, err := strconv.Atoi(c.QueryParam("limit")) // todo if zero?
	if err != nil {
		// todo
	}
	getPosts.Limit = int32(limit)

	since, errr := strconv.Atoi(c.QueryParam("since")) // todo if zero?
	if errr != nil {
		// todo
	}
	getPosts.Since = int64(since)

	sort := c.QueryParam("sort")
	getPosts.Sort = sort

	var desc bool
	desc, err = strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		// todo
	}
	getPosts.Desc = desc

	posts, errrr := h.ThreadsUcase.GetThreadPosts(slug, *getPosts)
	if errrr.Message == "404" {
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
