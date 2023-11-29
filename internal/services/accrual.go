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
	//defer wg.Done()
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

	wg.Done()
	return orderAccrual, nil
}

func (s *AccrualService) GetStatusOrder(outCh chan<- string, wg *sync.WaitGroup) {
	//outCh := make(chan string, 10)
	//defer wg.Done()
	//defer close(outCh)
	statusOrder, err := s.store.GetOrderStatus([]string{models.NEW})
	if err != nil {
		log.Errorf("failed get order by statused: %w", err)
	}
	//helpers.Infof("Sent:", statusOrder)
	//outCh <- statusOrder
	//statusOrder, err = s.store.GetOrderStatus([]string{models.PROCESSED)
	//if err != nil {
	//	log.Errorf("failed get order by statused: %w", err)
	//}
	//log.Print(statusOrder, "func GetStatusOrder")
	for _, order := range statusOrder {
		helpers.Infof("Sent:", order)
		outCh <- order
		//if len(outCh) > 0 {
		//	outCh <- order
		//}
	}

	//wg.Done()
	//	close(outCh)
	//	wg.Done()

	//log.Print(<-outCh)
	//return outCh
}

//func (s *AccrualService) UpdateOrderStatus() {
//	//var wg sync.WaitGroup
//	ordersChan := make(chan string, 10)
//	//defer close(ordersChan)
//	outCh := make(chan models.OrderAccrual, 10)
//	ticker := time.NewTicker(time.Duration(1) * time.Second)
//	defer ticker.Stop()
//	//wg.Add(2)
//	go func() {
//		s.GetStatusOrder(ordersChan)
//		for range ticker.C {
//			for {
//				select {
//				case <-ordersChan:
//					//close(ordersChan)
//					return
//				default:
//					helpers.Infof("get order number started")
//					//wg.Add(1)
//					s.GetStatusOrder(ordersChan)
//					//wg.Done()
//					helpers.Infof("get order number stoped")
//				}
//			}
//
//		}
//		//wg.Done()
//	}()
//	//wg.Add(1)
//	go func() {
//		for range ticker.C {
//			for {
//				select {
//				case <-outCh:
//					//	close(outCh)
//					return
//				default:
//					//if order, ok := <-ordersChan; ok {
//					//	if !ok {
//					//		break
//					//	}
//					for order := range ordersChan {
//						helpers.Infof("Received:", order)
//						resp, err := s.GetStatusAccrual(order,&wg)
//						if err != nil {
//							log.Infof("submit order: CalcOrderAccrual: %v", err)
//							return
//						}
//						helpers.Infof("Received started")
//						helpers.Infof(order, "order number")
//						log.Print(resp, "resp")
//						outCh <- resp
//					}
//
//				}
//			}
//		}
//		//wg.Done()
//	}()
//	//go func() {
//	//	wg.Wait()
//	//	defer close(ordersChan)
//	//	//defer close(outCh)
//	//}()
//	//wg.Wait()
//	go func() {
//		for {
//			select {
//			case <-outCh:
//				close(ordersChan)
//				//close(outCh)
//				return
//			default:
//				for order := range outCh {
//					helpers.Infof("START UPDATER")
//					log.Infof(order.Order, "responce order")
//					if order.Status == models.PROCESSED || order.Status == models.INVALID {
//						log.Print(order.Status, "order.status")
//						log.Print(order.Accrual, "accrual for begin update order")
//						err := s.store.UpdateOrder(order.Order, order.Status, order.Accrual)
//						if err != nil {
//							log.Print("error update order", err)
//							return
//						}
//						log.Info(order.Accrual, "accrual after update")
//						//close(a.OrdersChan)
//
//					}
//					if order.Status == models.PROCESSING || order.Status == models.REGISTERED {
//						log.Print(order.Status, "order.status")
//						err := s.store.UpdateOrderStatus(order.Order, order.Status)
//						if err != nil {
//							log.Print("error update order", err)
//							return
//						}
//						//close(a.OrdersChan)
//
//					}
//				}
//
//			}
//		}
//	}()
//	//	go startTask(ordersChan)
//
//}

func (s *AccrualService) CheckOrderStatus() {
	//wp := workerpool.New(3)
	var wg sync.WaitGroup
	ordersChan := make(chan string, 10)
	//defer close(ordersChan)
	outCh := make(chan models.OrderAccrual, 10)
	//defer close(outCh)
	//ticker := time.NewTicker(5 * time.Second)
	//defer ticker.Stop()
	pollTicker := time.NewTicker(time.Duration(1) * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Duration(1) * time.Second)
	defer reportTicker.Stop()
	//wg.Add(1)
	//go s.GetStatusOrder(ordersChan, &wg)
	//go func() {
	//	helpers.Infof("get order number started")
	//	for range pollTicker.C {
	//		statusOrder, err := s.store.GetOrderStatus([]string{models.NEW})
	//		if err != nil {
	//			log.Errorf("failed get order by statused: %w", err)
	//		}
	//		//statusOrder, err = s.store.GetOrderStatus([]string{models.PROCESSED)
	//		//if err != nil {
	//		//	log.Errorf("failed get order by statused: %w", err)
	//		//}
	//		log.Print(statusOrder, "func GetStatusOrder")
	//		for _, order := range statusOrder {
	//			helpers.Infof("Sent:", order)
	//			ordersChan <- order
	//		}
	//		//s.GetStatusOrder(ordersChan)
	//		helpers.Infof("get order number stoped")
	//	}
	//	//close(ordersChan)
	//	//wg.Done()
	//	//wp.StopWait()
	//	//wp.StopWait()
	//
	//}()

	//ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	//defer cancel()

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	//defer cancel ()
	//err := uh.storage.UpdateOrderStatus(order.ID, models.PROCESSING)
	//if err != nil {
	//	log.Println("update order status: %v", err)
	//}
	//wg.Add(1)
	go func() {
		//for range pollTicker.C {
		//statusOrder, err := s.store.GetOrderStatus([]string{models.NEW})
		//if err != nil {
		//	log.Errorf("failed get order by statused: %w", err)
		//}
		////statusOrder, err = s.store.GetOrderStatus([]string{models.PROCESSED)
		////if err != nil {
		////	log.Errorf("failed get order by statused: %w", err)
		////}
		//log.Print(statusOrder, "func GetStatusOrder")
		//for _, order := range statusOrder {
		//	helpers.Infof("Sent:", order)
		//	ordersChan <- order
		//}
		s.GetStatusOrder(ordersChan, &wg)
		helpers.Infof("get order number stoped")
		//close(ordersChan)
		//	}
		//	wg.Done()
	}()

	wg.Add(1)
	go func() {
		//defer close(ordersChan)
		helpers.Infof("Received started")
		for range reportTicker.C {
			//if order, ok := <-ordersChan; ok {
			for order := range ordersChan {
				helpers.Infof("Received:", order)
				resp, err := s.GetStatusAccrual(order, &wg)
				if err != nil {
					log.Infof("submit order: CalcOrderAccrual: %v", err)
					return
				}
				//
				//if resp.Status != string(models.PROCESSED) {
				//	time.AfterFunc(2*time.Second, func() {
				//		ordersChan <- order
				//	})
				//	return
				//}
				helpers.Infof(order, "order number")
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
	//go startTask(ordersChan)
	//go startTask1(outCh)

	//wg.Wait()
	go func() {
		for order := range outCh {
			helpers.Infof("START UPDATER")
			log.Infof(order.Order, "responce order")
			if order.Status == models.PROCESSED || order.Status == models.INVALID {
				log.Print(order.Status, "order.status")
				log.Print(order.Accrual, "accrual for begin update order")
				err := s.store.UpdateOrder(order.Order, order.Status, order.Accrual)
				if err != nil {
					log.Print("error update order", err)
					return
				}
				log.Info(order.Accrual, "accrual after update")
				//close(a.OrdersChan)

			}
			if order.Status == models.PROCESSING || order.Status == models.REGISTERED {
				log.Print(order.Status, "order.status")
				err := s.store.UpdateOrderStatus(order.Order, order.Status)
				if err != nil {
					log.Print("error update order", err)
					return
				}
				//close(a.OrdersChan)

			}
		}
		close(ordersChan)
		close(outCh)
		wg.Done()
	}()
	wg.Wait()
	//	close(outCh)
}

func (s *AccrualService) UpdateOrdersStatus() {
	//wp := workerpool.New(3)
	var wg sync.WaitGroup
	ordersChan := make(chan string, 10)
	//defer close(ordersChan)
	//outCh := make(chan models.OrderAccrual, 10)
	//defer close(outCh)
	//ticker := time.NewTicker(5 * time.Second)
	//defer ticker.Stop()
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
		//close(ordersChan)
		wg.Done()
	}()

	//defer close(ordersChan)

	//}

	wg.Add(1)
	//for i := 0; i < len(ordersChan); i++ {
	go func() {
		//defer close(ordersChan)
		helpers.Infof("Received started")
		for range reportTicker.C {
			//if order, ok := <-ordersChan; ok {
			for order := range ordersChan {
				helpers.Infof("Received:", order)
				resp, err := s.GetStatusAccrual(order, &wg)
				if err != nil {
					log.Infof("submit order: CalcOrderAccrual: %v", err)
					return
				}

				if resp.Status != models.PROCESSED {
					time.AfterFunc(3*time.Second, func() {
						ordersChan <- order
					})
					return
				}
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
					//close(a.OrdersChan)

				}

			}
		}
	}()
	//defer close(ordersChan)
	//}
	go startTask(ordersChan)
	wg.Wait()
	//	close(outCh)
}

//	wp.StopWait()
//	wg.Done()
//
//}()
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

//wg.Add(1)

//defer close(outCh)

//go func() {
//	//if order, ok := <-outCh; ok {
//	for order := range outCh {
//		log.Print(order.Order, "resp.order")
//		if order.Status == models.PROCESSED || order.Status == models.INVALID {
//			log.Print(order.Status, "order.status")
//			log.Print(order.Accrual, "accrual for begin update order")
//			err := s.store.UpdateOrder(order.Order, order.Status, order.Accrual)
//			if err != nil {
//				log.Print("error update order", err)
//				return
//			}
//			log.Print(order.Accrual, "accrual after update")
//			//close(a.OrdersChan)
//
//		}
//		if order.Status == models.PROCESSING || order.Status == models.REGISTERED {
//			log.Print(order.Status, "order.status")
//			err := s.store.UpdateOrderStatus(order.Order, order.Status)
//			if err != nil {
//				log.Print("error update order", err)
//				return
//			}
//			//close(a.OrdersChan)
//
//		}
//	}
//	//	}
//	//wg.Done()
//}()
//	return outCh
//defer close(ordersChan)
//defer close(outCh)
//}

func startTask(taskChan chan string) {
	for {
		select {
		case <-taskChan:
			//close(taskChan)
			return
		default:
			break
			//	log.Println("Задача выполняется...")
		}
	}
}

func startTask1(taskChan chan models.OrderAccrual) {
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
