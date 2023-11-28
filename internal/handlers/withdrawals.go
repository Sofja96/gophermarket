package handlers

import (
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetWithdrawals(storage *pg.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(models.ContextKeyUser).(string)
		withdrawals, err := storage.Getwithdrawals(user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Something went wrong")
		}
		if len(withdrawals) == 0 {
			return c.JSON(http.StatusNoContent, []models.Order{})
		}
		c.Response().Header().Set("Content-Type", "application/json")
		return c.JSON(http.StatusOK, withdrawals)

	}
}
