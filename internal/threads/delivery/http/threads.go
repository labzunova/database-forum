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
	slugOrId := c.Param("slug_or_id")
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		return err
	}

	thread, err :=  h.ThreadsUcase.UpdateThread(slugOrId, *newThread)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)
}

// ThreadGetPosts Получение списка сообщений в данной ветке форуме.
// Сообщения выводятся отсортированные по дате создания.
func (h *Handler) ThreadGetPosts(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	getPosts := new(models.ParseParamsThread)

	var err error
	getPosts.Limit, err = strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "limit не найден")
	}

	getPosts.Since, err = strconv.Atoi(c.QueryParam("since"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "since не найден")
	}

	sort := c.QueryParam("sort")
	getPosts.Sort = sort

	getPosts.Desc, err = strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "desc не найден")
	}

	posts, errr := h.ThreadsUcase.GetThreadPosts(slugOrId, *getPosts)
	if errr.Code == 404 {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, posts)
}

// ThreadVote Изменение голоса за ветвь обсуждения.
// Один пользователь учитывается только один раз и может изменить своё мнение.
func (h *Handler) ThreadVote(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newVote := new(models.Vote)
	if err := c.Bind(newVote); err != nil {
		return err
	}

	thread, err :=  h.ThreadsUcase.VoteThread(slug, *newVote)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Ветка отсутствует в форуме")
	}

	return c.JSON(http.StatusOK, thread)
}
