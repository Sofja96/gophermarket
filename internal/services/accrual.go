package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/labstack/gommon/log"
	"net/http"
)

type AccrualService struct {
	addr string
}

func NewAccrualService(addr string) *AccrualService {
	return &AccrualService{addr: addr}
}

func (s *AccrualService) GetStatusAccrual(ctx context.Context, orderNumber string) (models.OrderAccrual, error) {
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
	if resp.StatusCode != http.StatusOK {
		return orderAccrual, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
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

	return orderAccrual, nil
}
