package http

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"encoding/json"
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

	fullPost := posts.FullPost{}

	post, err := h.PostsUcase.GetPost(int(id))
	fmt.Println("get post error(handler)", err)
	if err.Code == 404 {
		err.Message = fmt.Sprintf("Can't find post with id: %d", id)
		return c.JSON(http.StatusNotFound, err)
	}

	if len(related) != 0 {
		fmt.Println("related exist")
		fullPost, err = h.PostsUcase.GetPostInfo(int(id), relatedSlice)
		fmt.Println("get post info error(handler)", err)
		if err.Code != 200 {
			return c.JSON(http.StatusInternalServerError, err.Code)
		}
	}

	fullPost.Post = &post
	fmt.Println("done")
	fmt.Println("post", fullPost)

	return c.JSON(http.StatusOK, fullPost)
}

// PostUpdate Изменение сообщения на форуме.
// Если сообщение поменяло текст, то оно должно получить отметку `isEdited`.
func (h *Handler) PostUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 0, 64)
	newMessage := new(posts.UpdateMessage)
	if err := c.Bind(newMessage); err != nil {
		return err
	}

	post, err := h.PostsUcase.UpdatePost(int(id), newMessage.Message)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, post)
}

// PostsCreate Добавление новых постов в ветку обсуждения на форум.
// Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).
func (h *Handler) PostsCreate(c echo.Context) error {
	slug := c.Param("slug_or_id")
	newPosts := make([]models.Post, 0)
	err := json.NewDecoder(c.Request().Body).Decode(&newPosts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	posts, errr := h.PostsUcase.CreatePosts(slug, newPosts)
	switch errr.Code {
	case 404:
		fmt.Println("not found")
		return c.JSON(http.StatusNotFound, errr)
	case 409:
		errr.Message = "Parent post was created in another thread"
		return c.JSON(http.StatusConflict, errr)
	}

	if len(newPosts) == 0 {
		return c.JSON(http.StatusCreated, newPosts)
	}

	return c.JSON(http.StatusCreated, posts)
}
