package http

import (
	"DBproject/helpers"
	"DBproject/internal/user"
	"DBproject/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	UserUcase user.UserUsecase
}

func NewUserHandler(userUcase user.UserUsecase) user.UserHandler {
	handler := &Handler{
		UserUcase: userUcase,
	}
	return handler
}

// UserCreate Создание нового пользователя в базе данных.
func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nickname := params["nickname"]

	newUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&newUser)
	fmt.Println(newUser)

	newUser.Nickname = nickname

	err := h.UserUcase.Create(*newUser)

	// если такой уже есть
	if err.Code == 409 {
		users, errr := h.UserUcase.GetExistingUsers(newUser.Nickname, newUser.Email)
		if errr.Code != 200 {
			helpers.CreateResponse(w, http.StatusInternalServerError, "error")
			return
		}
		helpers.CreateResponse(w, http.StatusConflict, users)
		return
	}
	fmt.Println(newUser, "createsd")

	helpers.CreateResponse(w, http.StatusCreated, newUser)
	return
}

// UserGetOne Получение информации о пользователе форума по его имени.
func (h *Handler) UserGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nickname := params["nickname"]

	user, err := h.UserUcase.GetByNickname(nickname)
	if err.Code == 404 {
		err.Message = "Can't find user by nickname: " + nickname
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
	}

	helpers.CreateResponse(w, http.StatusOK, user)
	return
}

// UserUpdate Изменение информации в профиле пользователя.
func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	nickname := params["nickname"]
	newUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&newUser)

	newUser.Nickname = nickname

	user, err := h.UserUcase.Update(*newUser)
	fmt.Println(err)
	switch err.Code {
	case 404:
		err.Message = "Can't find user by nickname: " + nickname
		helpers.CreateResponse(w, http.StatusNotFound, err)
		return
		case 409:
		users, errr := h.UserUcase.GetExistingUsers(nickname, newUser.Email)
		errr.Message = "This email is already registered by user: " + users[0].Nickname
			helpers.CreateResponse(w, http.StatusConflict, errr)
			return
	}

	helpers.CreateResponse(w, http.StatusOK, user)
	return
}
