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
	postsRepository posts.PostsRepo
}

func NewPostsHandler(repo posts.PostsRepo) posts.PostsHandler {
	handler := &Handler{
		postsRepository: repo,
	}

	return handler
}

// PostGetOne Получение информации о ветке обсуждения по его имени.
func (h *Handler) PostGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 0, 64)
	var related []string
	related = strings.Split(r.URL.Query().Get("related"), ",")

	fullPost := posts.FullPost{}

	post, err := h.postsRepository.GetPost(int(id))
	if err.Code == 404 {
		err.Message = fmt.Sprintf("Can't find post with id: %d", id)
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	if len(related) != 0 {
		for _, info := range related {
			switch info {
			case "user":
				user, _ := h.postsRepository.GetPostAuthor(post.Author)
				fmt.Println("post author", user)
				fullPost.User = &user
			case "forum":
				forum, _ := h.postsRepository.GetPostForum(post.Forum)
				fmt.Println("post forum", forum)
				fullPost.Forum = &forum
			case "thread":
				thread, _ := h.postsRepository.GetPostThread(post.Thread)
				fmt.Println("post thread", thread)
				fullPost.Thread = &thread
			}
		}
	}
	fullPost.Post = &post

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

	post, err := h.postsRepository.UpdatePost(int(id), newMessage.Message)
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

	var thread models.Thread
	id, errID := strconv.Atoi(slug)
	if errID == nil {
		var errr models.Error
		thread, errr = h.postsRepository.GetThreadAndForumById(id)
		if errr.Code == 404 {
			errr.Message = "Can't find post thread by id: " + slug
			helpers.CreateResponse(w, http.StatusNotFound, err)
			return
		}
	} else {
		var errr models.Error
		thread, errr = h.postsRepository.GetThreadAndForumBySlug(slug)
		fmt.Println(errr)
		if errr.Code == 404 {
			errr.Message = "Can't find post thread by id: " + slug
			helpers.CreateResponse(w, http.StatusNotFound, err)
			return
		}
	}

	posts, errr := h.postsRepository.CreatePosts(thread, newPosts)
	if errr.Code == 404 {
		errr.Message = "Can't find user "
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}
	if errr.Code == 409 {
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
