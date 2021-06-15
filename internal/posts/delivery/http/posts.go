package http

import (
	"DBproject/helpers"
	"DBproject/internal/posts"
	"DBproject/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

// PostGetOne Получение информации о ветке обсуждения по его имени.
func (h *Handler) PostGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 0, 64)
	var related []string
	related = strings.Split(r.URL.Query().Get("related"), ",")
fmt.Println(related)
	fmt.Println("GET POST", id)

	fullPost := posts.FullPost{}

	post, err := h.PostsUcase.GetPost(int(id))
	fmt.Println("get post error(handler)", err)
	if err.Code == 404 {
		err.Message = fmt.Sprintf("Can't find post with id: %d", id)
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	if len(related) != 0 {
		fmt.Println("related exist")
		fullPost, err = h.PostsUcase.GetPostInfo(post, int(id), related)
		fmt.Println("get post info error(handler)", err)
		if err.Code != 200 {
			helpers.CreateResponse(w, http.StatusInternalServerError, err.Code)
			return
		}
	}

	fullPost.Post = &post
	fmt.Println("done")
	fmt.Println("post", fullPost)


	helpers.CreateResponse(w, http.StatusOK, fullPost)
	return
}

// PostUpdate Изменение сообщения на форуме.
// Если сообщение поменяло текст, то оно должно получить отметку `isEdited`.
func (h *Handler) PostUpdate(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 0, 64)
	newMessage := new(posts.UpdateMessage)
	json.NewDecoder(r.Body).Decode(&newMessage)

	post, err := h.PostsUcase.UpdatePost(int(id), newMessage.Message)
	if err.Code == 404 {
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	helpers.CreateResponse(w, http.StatusOK, post)
	return
}

// PostsCreate Добавление новых постов в ветку обсуждения на форум.
// Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).
func (h *Handler) PostsCreate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug_or_id"]

	newPosts := make([]models.Post, 0)
	err := json.NewDecoder(r.Body).Decode(&newPosts)
	if err != nil {
		helpers.CreateResponse(w, http.StatusBadRequest, err)
		return
	}

	posts, errr := h.PostsUcase.CreatePosts(slug, newPosts)
	switch errr.Code {
	case 404:
		fmt.Println("not found")
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	case 409:
		errr.Message = "Parent post was created in another thread"
		helpers.CreateResponse(w, http.StatusConflict, errr)
		return
	}

	if len(newPosts) == 0 {
		helpers.CreateResponse(w, http.StatusCreated, newPosts)
		return
	}

	helpers.CreateResponse(w, http.StatusCreated, posts)
	return
}
