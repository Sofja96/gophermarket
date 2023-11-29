package services

import (
	"encoding/json"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"strconv"
	"sync"

	//"github.com/gammazero/workerpool"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type AccrualService struct {
	addr  string
	store *pg.Postgres
}

func NewAccrualService(addr string, storage *pg.Postgres) *AccrualService {
	return &AccrualService{addr: addr, store: storage}
}

func (s *AccrualService) GetStatusAccrual(orderNumber string, wg *sync.WaitGroup) (models.OrderAccrual, error) {

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

func (s *AccrualService) GetStatusOrder(outCh chan<- string, wg *sync.WaitGroup) {
	statusOrderNew, err := s.store.GetOrderStatus([]string{models.NEW})
	if err != nil {
		log.Errorf("failed get order by statused: %w", err)
	}
	statusOrderProc, err := s.store.GetOrderStatus([]string{models.PROCESSING})
	if err != nil {
		log.Errorf("failed get order by statused: %w", err)
	}

	for _, order := range statusOrderNew {
		helpers.Infof("Sent New:", order)
		outCh <- order
	}
	for _, order := range statusOrderProc {
		helpers.Infof("Sent Proc:", order)
		outCh <- order
	}
}

func (s *AccrualService) UpdateOrdersStatus() {
	var wg sync.WaitGroup
	ordersChan := make(chan string, 10)
	pollTicker := time.NewTicker(time.Duration(1) * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Duration(1) * time.Second)
	defer reportTicker.Stop()

	s.GetStatusOrder(ordersChan, &wg)
	wg.Add(1)
	go func() {
		helpers.Infof("get order number stoped")
		for range pollTicker.C {
			s.GetStatusOrder(ordersChan, &wg)
			log.Print(len(ordersChan), "lenth channel")
			helpers.Infof(strconv.Itoa(cap(ordersChan)), "cap channel")
			helpers.Infof("get order number stoped")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		helpers.Infof("Received started")
		for range reportTicker.C {
			for order := range ordersChan {
				helpers.Infof("Received:", order)
				resp, err := s.GetStatusAccrual(order, &wg)
				if err != nil {
					log.Infof("submit order: CalcOrderAccrual: %v", err)
					return
				}

				//if resp.Status != models.PROCESSED {
				//	time.AfterFunc(3*time.Second, func() {
				//		ordersChan <- order
				//	})
				//	return
				//}
				helpers.Infof(order, "order number")
				log.Print(resp, "resp")
				helpers.Infof("START UPDATER")
				log.Infof(resp.Order, "responce order")
				if resp.Status == models.PROCESSED || resp.Status == models.INVALID {
					log.Print(resp.Status, "order.status")
					log.Print(resp.Accrual, "accrual for begin update order")
					err := s.store.UpdateOrder(resp.Order, resp.Status, resp.Accrual)
					if err != nil {
						log.Print("error update order", err)
						return
					}
					log.Info(resp.Accrual, "accrual after update")
				}

			}
		}
	}()
	//go startTask(ordersChan)
	wg.Wait()
}

func startTask(taskChan chan string) {
	for {
		select {
		case <-taskChan:
			return
		default:
			break
			//	log.Println("Задача выполняется...")
		}
	}
}
