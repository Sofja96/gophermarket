package handlers

import (
	_ "context"
	"errors"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/services"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"sync"
	_ "time"
)

type UserHandler struct {
	storage        *pg.Postgres
	accrualService *services.AccrualService
	//pool           submitter
}

func NewUserHandler(storage *pg.Postgres, as *services.AccrualService) *UserHandler {
	return &UserHandler{storage: storage, accrualService: as}
}

// канал для отправки данных о номере заказа
func PostOrder(storage *pg.Postgres, ordersChan chan<- string, wg *sync.WaitGroup) echo.HandlerFunc {
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
		log.Println("SubmitOrder with order number: %s", orderNumber)
		if len(orderNumber) == 0 {
			return c.String(http.StatusBadRequest, "Empty request body")
		}

		if isValid := helpers.IsOrderNumberValid(orderNumber); !isValid {
			return c.String(http.StatusUnprocessableEntity, "Wrong number of order format")
		}
		username := c.Get(models.ContextKeyUser).(string)
		log.Println(c.Get(models.ContextKeyUser))

		//	order, err := uh.storage.CreateOrder(orderNumber, username)
		_, err = storage.CreateOrder(orderNumber, username)
		if err != nil {
			//switch {
			//case errors.Is(err, helpers.ErrAnotherUserOrder):
			//	log.Println(err)
			//	return c.String(http.StatusConflict, "order number already exists for another user")
			//case errors.Is(err, helpers.ErrExistsOrder):
			//	log.Println(err)
			//	return c.String(http.StatusOK, "order number already exists")
			////errors.New()
			//default:
			//	return c.String(http.StatusInternalServerError, "Something went wrong")
			//
			//}
			if errors.Is(err, helpers.ErrAnotherUserOrder) {
				return c.String(http.StatusConflict, "order number already exists for another user")
			}
			if errors.Is(err, helpers.ErrExistsOrder) {
				return c.String(http.StatusOK, "order number already exists")
			}
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
		ordersChan <- orderNumber
		//val := <-ordersChan
		//println(val)
		log.Println(orderNumber, username, "order+user in create")
		//	uh.calcAndApplyAccrual(order, username)
		//wg.Done() // decrement counter
		return c.String(http.StatusAccepted, "")

	}
}

func (uh *UserHandler) calcAndApplyAccrual(order *models.Order, userID string) {
	go func() {
		//ticker := time.NewTicker(1 * time.Second)
		//defer ticker.Stop()
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		//defer cancel()
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		//defer cancel ()
		//err := uh.storage.UpdateOrderStatus(order.ID, models.PROCESSING)
		//if err != nil {
		//	log.Println("update order status: %v", err)
		//}
		//	for range ticker.C {
		resp, err := uh.accrualService.GetStatusAccrual(order.Number)
		if err != nil {
			log.Println("submit order: CalcOrderAccrual: %v", err)
			return
		}
		log.Println(resp, "resp")
		//err = uh.storage.UpdateOrderAccrualAndUserBalance(ctx, order.ID, userID, resp)
		//if err != nil {
		//	log.Println("submit order: UpdateOrderAccrualAndUserBalance: %v", err)
		//	return
		//}
		//err = uh.storage.UpdateOrderStatus(ctx, order.ID, models.PROCESSED)
		//if err != nil {
		//	log.Println("update order status: %v", err)
		//}
		//	}
	}()
}
