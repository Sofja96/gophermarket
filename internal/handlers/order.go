package handlers

import (
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/services"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	storage        *pg.Postgres
	accrualService *services.AccrualService
	//pool           submitter
}

func NewUserHandler(storage *pg.Postgres, as *services.AccrualService) *UserHandler {
	return &UserHandler{storage: storage, accrualService: as}
}

func PostOrder(uh *UserHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "text/plain" {
			return c.String(http.StatusUnsupportedMediaType, "")
		}
		username, _, ok := c.Request().BasicAuth()
		if !ok {
			return c.String(http.StatusInternalServerError, "")
		}
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to read request body")
		}
		defer c.Request().Body.Close()

		orderNumber := string(body)
		log.Println("SubmitOrder with order number: %s", orderNumber)
		if len(orderNumber) == 0 {
			return c.String(http.StatusBadRequest, "Empty request body")
		}

		if isValid := helpers.IsOrderNumberValid(orderNumber); !isValid {
			return c.String(http.StatusUnprocessableEntity, "Wrong number of order format")
		}

		order, err := uh.storage.CreateOrder(orderNumber, username)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
		uh.calcAndApplyAccrual(order, username)
		return c.String(http.StatusAccepted, "")

	}
}

func (uh *UserHandler) calcAndApplyAccrual(order *models.Order, userID string) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		//defer cancel ()
		err := uh.storage.UpdateOrderStatus(order.ID, models.PROCESSING)
		if err != nil {
			log.Println("update order status: %v", err)
		}
		for range ticker.C {
			resp, err := uh.accrualService.GetStatusAccrual(order.Number)
			if err != nil {
				log.Println("submit order: CalcOrderAccrual: %v", err)
				return
			}
			err = uh.storage.UpdateOrderAccrualAndUserBalance(order.ID, userID, resp)
			if err != nil {
				log.Println("submit order: UpdateOrderAccrualAndUserBalance: %v", err)
				return
			}
			err = uh.storage.UpdateOrderStatus(order.ID, models.PROCESSED)
			if err != nil {
				log.Println("update order status: %v", err)
			}
		}
	}()
}
