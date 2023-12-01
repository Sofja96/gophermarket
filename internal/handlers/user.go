package handlers

import (
	"encoding/json"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/middleware"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterUser(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		var newUser models.User
		var err error
		if err := json.NewDecoder(c.Request().Body).Decode(&newUser); err != nil {
			return c.JSON(http.StatusBadRequest, "")
		}
		if len(newUser.Login) == 0 && len(newUser.Password) == 0 {
			return c.JSON(http.StatusBadRequest, "empty credentials")
		}
		existingUser, err := storage.GetUserIDByName(newUser.Login)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Something went wrong")
		}
		if existingUser {
			return c.JSON(http.StatusConflict, "Login is already exists")
		}

		hash, err := helpers.HashPassword(newUser.Password)
		if err != nil {
			return err
		}
		_, err = storage.CreateUser(newUser.Login, hash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "user creation error")
		}

		token, err := middleware.CreateToken(newUser.Login)
		if err != nil {
			helpers.Error("error create token")
			return err
		}
		var bearer = "Bearer " + token
		c.Response().Header().Set("Authorization", bearer)
		return c.JSON(http.StatusOK, newUser)

	}
}

func LoginUser(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		var newUser models.User
		var err error

		if err := json.NewDecoder(c.Request().Body).Decode(&newUser); err != nil {
			return c.JSON(http.StatusBadRequest, "")
		}
		if len(newUser.Login) == 0 && len(newUser.Password) == 0 {
			return c.JSON(http.StatusBadRequest, "empty credentials")
		}
		existingUser, err := storage.GetUserIDByName(newUser.Login)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Something went wrong")
		}
		if !existingUser {
			return c.JSON(http.StatusConflict, "Users not found, please to registration")
		}
		hash, err := storage.GetUserHashPassword(newUser.Login)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Something went wrong")
		}
		err = helpers.CheckPassword(newUser.Password, hash)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "password is not coorect")
		}
		token, err := middleware.CreateToken(newUser.Login)
		if err != nil {
			helpers.Error("error create token")
			return err
		}
		var bearer = "Bearer " + token

		c.Response().Header().Set("Authorization", bearer)
		return c.JSON(http.StatusOK, newUser)

	}
}
