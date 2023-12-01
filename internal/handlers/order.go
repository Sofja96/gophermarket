package handlers

import (
	_ "context"
	"errors"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	_ "time"
)

func PostOrder(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "text/plain" {
			return c.String(http.StatusUnsupportedMediaType, "")
		}

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to read request body")
		}
		defer c.Request().Body.Close()

		orderNumber := string(body)
		helpers.Infof("SubmitOrder with order number: %s\n", orderNumber)
		if len(orderNumber) == 0 {
			return c.String(http.StatusBadRequest, "Empty request body")
		}

		if isValid := helpers.IsOrderNumberValid(orderNumber); !isValid {
			return c.String(http.StatusUnprocessableEntity, "Wrong number of order format")
		}
		username := c.Get(models.ContextKeyUser).(string)

		_, err = storage.CreateOrder(orderNumber, username)
		if err != nil {
			if errors.Is(err, helpers.ErrAnotherUserOrder) {
				return c.String(http.StatusConflict, "order number already exists for another user")
			}
			if errors.Is(err, helpers.ErrExistsOrder) {
				return c.String(http.StatusOK, "order number already exists")
			}
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
		return c.String(http.StatusAccepted, "")

	}
}

func GetOrders(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(models.ContextKeyUser).(string)
		orders, err := storage.GetOrders(user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Something went wrong")
		}
		if len(orders) == 0 {
			return c.JSON(http.StatusNoContent, []models.Order{})
		}
		c.Response().Header().Set("Content-Type", "application/json")
		return c.JSON(http.StatusOK, orders)

	}
}

func GetBalance(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(models.ContextKeyUser).(string)
		balance, err := storage.GetBalance(user)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}

		c.Response().Header().Set("Content-Type", "application/json")
		return c.JSON(http.StatusOK, balance)
	}
}
