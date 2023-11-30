package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
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

	helpers.Infof("CalcOrderAccrual response status: %d", resp.StatusCode)

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
		helpers.Infof("CalcOrderAccrual response: %s", r)
	}

	return orderAccrual, nil
}

func (s *AccrualService) GetStatusOrder(outCh chan<- string) {
	statusOrderNew, err := s.store.GetOrderStatus([]string{models.NEW})
	if err != nil {
		helpers.Error("failed get order by status: %s", err)
	}
	statusOrderProc, err := s.store.GetOrderStatus([]string{models.PROCESSING})
	if err != nil {
		helpers.Error("failed get order by status: %s", err)
	}

	for _, order := range statusOrderNew {
		outCh <- order
	}
	for _, order := range statusOrderProc {
		outCh <- order
	}
}

func (s *AccrualService) UpdateOrdersStatus(ctx context.Context) {
	var wg sync.WaitGroup
	ordersChan := make(chan string, 10)
	pollTicker := time.NewTicker(time.Duration(1) * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Duration(1) * time.Second)
	defer reportTicker.Stop()

	s.GetStatusOrder(ordersChan)
	wg.Add(1)
	go func() {
		for range pollTicker.C {
			select {
			case <-ctx.Done():
				close(ordersChan)
				return
			default:
				s.GetStatusOrder(ordersChan)
			}
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for range reportTicker.C {
			select {
			case <-ctx.Done():
				close(ordersChan)
				return
			default:
				for order := range ordersChan {
					resp, err := s.GetStatusAccrual(order, &wg)
					if err != nil {
						helpers.Error("submit order: CalcOrderAccrual: %v", err)
						return
					}
					if resp.Status == models.PROCESSED || resp.Status == models.INVALID {
						err := s.store.UpdateOrder(resp.Order, resp.Status, resp.Accrual)
						if err != nil {
							helpers.Error("error update OrderAccrual: %s", err)
							return
						}
					}

				}
			}
		}
	}()
	wg.Wait()
}
