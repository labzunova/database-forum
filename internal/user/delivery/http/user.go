package http

import (
	"DBproject/internal/user"
	"DBproject/models"
	"github.com/labstack/echo/v4"
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
func (h *Handler) UserCreate(c echo.Context) error {
	nickname := c.Param("nickname")
	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		return err
	}
	newUser.Nickname = nickname

	err := h.UserUcase.Create(*newUser)

	// если такой уже есть
	if err.Code == 409 {
		users, errr := h.UserUcase.GetExistingUsers(newUser.Nickname, newUser.Email)
		if errr.Code != 200 {
			return c.JSON(http.StatusInternalServerError, "error")
		}
		return c.JSON(http.StatusConflict, users)
	}

	return c.JSON(http.StatusOK, newUser)
}

// UserGetOne Получение информации о пользователе форума по его имени.
func (h *Handler) UserGetOne(c echo.Context) error {
	nickname := c.Param("nickname")

	user, err := h.UserUcase.GetByNickname(nickname)
	if err.Code == 404 {
		return c.JSON(http.StatusNotFound, "Пользователь отсутсвует в системе")
	}

	return c.JSON(http.StatusOK, user)
}

// UserUpdate Изменение информации в профиле пользователя.
func (h *Handler) UserUpdate(c echo.Context) error {
	nickname := c.Param("nickname")
	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		return err
	}
	newUser.Nickname = nickname

	user, err := h.UserUcase.Update(*newUser)
	switch err.Code {
	case 404:
		return c.JSON(http.StatusNotFound, "Пользователь отсутсвует в системе")
	case 409:
		return c.JSON(http.StatusConflict, "Новые данные профиля пользователя конфликтуют с имеющимися пользователями")
	}

	return c.JSON(http.StatusOK, user)
}
