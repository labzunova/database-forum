package http

import (
	"DBproject/helpers"
	"DBproject/internal/forum"
	"DBproject/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
func (h Handler) ForumCreate(w http.ResponseWriter, r *http.Request) {
	newForum := new(models.Forum)
	json.NewDecoder(r.Body).Decode(&newForum)

	forum, errr := h.ForumUcase.CreateNewForum(*newForum)
	switch errr.Code {
	case 404:
		errr.Message = "Can't find user with nickname: " + newForum.User
		helpers.CreateResponse(w, http.StatusNotFound, errr)
		return
	case 409:
		forumOld, _ := h.ForumUcase.GetForum(newForum.Slug)
		helpers.CreateResponse(w, http.StatusConflict,forumOld)
		return
	}

	helpers.CreateResponse(w, http.StatusCreated, forum)
	return
}

// ForumGetOne Получение информации о форуме по его идентификаторе
func (h Handler) ForumGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	forumResponse, err := h.ForumUcase.GetForum(slug)
	if err.Code != 200 {
		//err.Message = "Can't find forum with slug: " + slug
		err.Message = "Can't find forum"
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}


	helpers.CreateResponse(w, http.StatusOK, forumResponse)
	return
}

// ThreadCreate Добавление новой ветки обсуждения на форум
func (h Handler) ThreadCreate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	newThread := new(models.Thread)
	json.NewDecoder(r.Body).Decode(&newThread)


	thread, err := h.ForumUcase.CreateThread(slug, *newThread)

	switch err.Code {
	case 404:
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	case 409:
		helpers.CreateResponse(w, http.StatusConflict, thread)
		return
	}

	helpers.CreateResponse(w, http.StatusCreated, thread)
	return
}

// ForumGetUsers Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
// Пользователи выводятся отсортированные по nickname в порядке возрастания.
// Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.
func (h Handler) ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	getUsers := new(models.ParseParams)

	var errParse error
	getUsers.Limit, errParse = strconv.Atoi(r.URL.Query().Get("limit"))
	if errParse != nil {
		getUsers.Limit = 100000
	}
	getUsers.Since = r.URL.Query().Get("since")
	getUsers.Desc, errParse = strconv.ParseBool(r.URL.Query().Get("desc"))
	if errParse != nil {
		getUsers.Desc = false
	}

	users, err := h.ForumUcase.GetUsers(slug, *getUsers)
	if err.Code == 404 {
		err.Message = "Can't find forum by slug: " + slug
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	helpers.CreateResponse(w, http.StatusOK, users)
	return
}

// ForumGetThreads Получение списка ветвей обсужления данного форума.
// Ветви обсуждения выводятся отсортированные по дате создания.
func (h Handler) ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	getThreads := new(models.ParseParams)


	var errParse error
	getThreads.Limit, errParse = strconv.Atoi(r.URL.Query().Get("limit"))
	if errParse != nil {
		getThreads.Limit = 100000
	}
	getThreads.Since = r.URL.Query().Get("since")
	getThreads.Desc, errParse = strconv.ParseBool(r.URL.Query().Get("desc"))
	if errParse != nil {
		getThreads.Desc = false
	}

	fmt.Println(getThreads)

	threads, err := h.ForumUcase.GetThreads(slug, *getThreads)
	if err.Code != 200 {
		err.Message = "Can't find forum by slug: " + slug
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return	}

	fmt.Println("ФЛОООООООО")
	helpers.CreateResponse(w, http.StatusOK, threads)
	return
}
