package services

import (
	"encoding/json"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
	"time"
)

type AccrualService struct {
	addr  string
	store *pg.Postgres
}

func NewAccrualService(addr string, storage *pg.Postgres) *AccrualService {
	return &AccrualService{addr: addr, store: storage}
}

func (s *AccrualService) GetStatusAccrual(orderNumber string) (models.OrderAccrual, error) {
	var orderAccrual models.OrderAccrual
	url := fmt.Sprintf("%s/api/orders/%s", s.addr, orderNumber)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return orderAccrual, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return orderAccrual, err
	}
	defer resp.Body.Close()

	log.Infof("CalcOrderAccrual response status: %d", resp.StatusCode)
	//if resp.StatusCode != http.StatusOK {
	//	return orderAccrual, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	//}
	if resp.StatusCode == http.StatusNoContent {
		return orderAccrual, fmt.Errorf("order not registered: %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return orderAccrual, fmt.Errorf("too many request: %d", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return orderAccrual, fmt.Errorf("internal accrual system error: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&orderAccrual)
	if err != nil {
		return orderAccrual, err
	}

	r, err := json.Marshal(orderAccrual)
	if err == nil {
		log.Infof("CalcOrderAccrual response: %s", r)
	}
	log.Print(resp, "check resp")

	return orderAccrual, nil
}

func (s *AccrualService) CheckOrderStatus(ordersChan chan string) {
	var wg sync.WaitGroup
	outCh := make(chan models.OrderAccrual, 10)
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		//ctx := context.Background()
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		//defer cancel()

		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		//defer cancel ()
		//err := uh.storage.UpdateOrderStatus(order.ID, models.PROCESSING)
		//if err != nil {
		//	log.Println("update order status: %v", err)
		//}
		for range ticker.C {
			for order := range ordersChan {
				resp, err := s.GetStatusAccrual(order)
				if err != nil {
					log.Infof("submit order: CalcOrderAccrual: %v", err)
					return
				}

				if resp.Status != string(models.PROCESSED) {
					time.AfterFunc(2*time.Second, func() {
						ordersChan <- order
					})
					return
				}
				log.Infof(order, "order number")
				log.Print(resp, "resp")
				outCh <- resp
				//if resp.Status != string(models.PROCESSED) {
				//	time.AfterFunc(5*time.Second, func() {
				//		//ordersChan <- order
				//	})
				//	return
				//}
				//log.Println(resp.Order, "resp.order")
				//if resp.Status == string(models.PROCESSED) {
				//	log.Println(order, "processed")
				//	//close(a.OrdersChan)
				//
				//}

				//	close(outCh)

				//close(a.OrdersChan)
				//close(a.OrdersChan)
				//	a.OrdersChan <- order
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

			}
		}
		wg.Done()

	}()
	//wg.Add(1)
	//go func() {
	//	ticker := time.NewTicker(time.Duration(1) * time.Second)
	//	defer ticker.Stop()
	//	newOrders, err := a.store.GetOrderStatus(models.NEW)
	//	if err != nil {
	//		if err != nil {
	//			log.Println("error get status order", err)
	//
	//		}
	//		for range ticker.C {
	//			//	var ordernum models.OrderAccrual
	//			//var ordernum string
	//			for _, ordernum := range newOrders {
	//				outCh <- ordernum
	//			}
	//		}
	//	}
	//	wg.Done()
	//}()
	wg.Add(1)
	go func() {
		for order := range outCh {
			log.Print(order.Order, "resp.order")
			if order.Status == string(models.PROCESSED) || order.Status == string(models.INVALID) {
				log.Print(order.Status, "order.status")
				log.Print(order.Accrual, "accrual for begin update order")
				err := s.store.UpdateOrder(order.Order, order.Status, order.Accrual)
				if err != nil {
					log.Print("error update order", err)
					return
				}
				log.Print(order.Accrual, "accrual after update")
				//close(a.OrdersChan)

			}
			if order.Status == string(models.PROCESSING) || order.Status == string(models.REGISTERED) {
				log.Print(order.Status, "order.status")
				err := s.store.UpdateOrderStatus(order.Order, models.OrderStatus(order.Status))
				if err != nil {
					log.Print("error update order", err)
					return
				}
				//close(a.OrdersChan)

			}
		}
		wg.Done()
	}()
	//	return outCh
	//go startTask(ordersChan)
	//wg.Wait()
}
