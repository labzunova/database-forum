package http

import (
	"DBproject/helpers"
	"DBproject/internal/threads"
	"DBproject/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	threadsRepo threads.ThreadsRepo

}

func NewThreadsHandler(repo threads.ThreadsRepo) threads.ThreadsHandler {
	handler := &Handler{
		threadsRepo: repo,
	}
	return handler
}

// ThreadGetOne Получение информации о ветке обсуждения по его имени.
func (h *Handler) ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug_or_id"]

	id := h.SlugOrID(slug)
	thread, err := h.threadsRepo.GetThread(slug, id)
	if err.Code == 404 {
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	helpers.CreateResponse(w, http.StatusOK, thread)
	return
}

// ThreadUpdate Обновление ветки обсуждения на форуме.
func (h *Handler) ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slugOrId := params["slug_or_id"]

	newThread := new(models.Thread)
	json.NewDecoder(r.Body).Decode(&newThread)

	var err models.Error
	id := h.SlugOrID(slugOrId)
	if id != 0 {
		*newThread, err = h.threadsRepo.UpdateThreadById(id, *newThread)
	}
	thread, err := h.threadsRepo.UpdateThreadBySlug(slugOrId, *newThread)
	//thread, err := h.ThreadsUcase.UpdateThread(slugOrId, *newThread)
	if err.Code == 404 {
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	helpers.CreateResponse(w, http.StatusOK, thread)
	return
}

// ThreadGetPosts Получение списка сообщений в данной ветке форуме.
// Сообщения выводятся отсортированные по дате создания.
func (h *Handler) ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slugOrId := params["slug_or_id"]
	getPosts := new(models.ParseParamsThread)

	var errParse error
	getPosts.Limit, errParse = strconv.Atoi(r.URL.Query().Get("limit"))
	if errParse != nil {
		getPosts.Limit = 100000
	}

	getPosts.Since, errParse = strconv.Atoi(r.URL.Query().Get("since"))
	if errParse != nil {
		getPosts.Since = 0
	}

	getPosts.Sort = r.URL.Query().Get("sort")
	getPosts.Desc, errParse = strconv.ParseBool(r.URL.Query().Get("desc"))
	if errParse != nil {
		getPosts.Desc = false
	}

	id := h.SlugOrID(slugOrId)
	var err models.Error
	if id != 0 {
		fmt.Println("id:", id)
		id, err = h.threadsRepo.GetThreadIDBySlug(slugOrId, id)
		if err.Code != 200 {
			err.Message = "Can't find forum by slug: " + slugOrId
			helpers.CreateResponse(w, http.StatusNotFound, err)
			return
		}
	} else {
		id, err = h.threadsRepo.GetThreadIDBySlug(slugOrId, 0)
		fmt.Println("get by slug")
		if err.Code != 200 {
			err.Message = "Can't find forum by slug: " + slugOrId
			helpers.CreateResponse(w, http.StatusNotFound, err)
			return
		}
	}
	posts, err := h.threadsRepo.GetThreadPostsById(id, slugOrId, *getPosts)
	if err.Code == 404 {
		err.Message = "Can't find forum by slug: " + slugOrId
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	helpers.CreateResponse(w, http.StatusOK, posts)
	return
}

// ThreadVote Изменение голоса за ветвь обсуждения.
// Один пользователь учитывается только один раз и может изменить своё мнение.
func (h *Handler) ThreadVote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("VOTE handler")
	newVote := new(models.Vote)
	json.NewDecoder(r.Body).Decode(&newVote)

	params := mux.Vars(r)
	slug := params["slug_or_id"]


	id := h.SlugOrID(slug)
	fmt.Println("slug or id ", slug, id)
	if id != 0 {
		fmt.Println("vote by id")
		//err := h.threadsRepo.VoteThreadById(id, *newVote)
		//if err.Code != 200 {
		//	return models.Thread{}, err
		//}
		_ = h.threadsRepo.VoteThreadById(id, *newVote)

		thread, err := h.threadsRepo.GetThread("", id)
		if err.Code == 404 {
			helpers.CreateResponse(w, http.StatusNotFound, err)
			return
		}
		if err.Code != 200 {
			helpers.CreateResponse(w, http.StatusInternalServerError, "Error")
			return
		}
		helpers.CreateResponse(w, http.StatusOK, thread)
		return
	}

	fmt.Println("vote by slug")
	_ = h.threadsRepo.VoteThreadBySlug(slug, *newVote)

	thread, err := h.threadsRepo.GetThread(slug, 0)
	if err.Code == 404 {
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}
	if err.Code != 200 {
		helpers.CreateResponse(w, http.StatusInternalServerError, "Error")
		return
	}

	helpers.CreateResponse(w, http.StatusOK, thread)
	return
}


func (h *Handler) SlugOrID(slugOrId string) int {
	id, errID := strconv.Atoi(slugOrId)
	if errID == nil {
		return id
	}

	return 0
}
