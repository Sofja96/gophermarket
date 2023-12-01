package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func WithdrawBalance(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		var withDrawal models.UserWithdrawal

		if err := json.NewDecoder(c.Request().Body).Decode(&withDrawal); err != nil {
			return c.JSON(http.StatusBadRequest, "")
		}

		if isValid := helpers.IsOrderNumberValid(withDrawal.Order); !isValid {
			return c.String(http.StatusUnprocessableEntity, "Wrong number of order format")
		}
		username := c.Get(models.ContextKeyUser).(string)

		err := storage.WithdrawBalance(username, withDrawal.Order, withDrawal.Sum)
		if err != nil {
			if errors.Is(err, helpers.ErrInsufficientBalance) {
				return c.String(http.StatusPaymentRequired, "insufficient balance")
			}
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}

		c.Response().Header().Set("Content-Type", "application/json")

		return c.JSON(http.StatusOK, withDrawal)

	}
}
