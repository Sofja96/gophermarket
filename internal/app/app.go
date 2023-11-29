package app

import (
	"context"
	"github.com/Sofja96/gophermarket.git/internal/app/config"
	"github.com/Sofja96/gophermarket.git/internal/handlers"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"github.com/Sofja96/gophermarket.git/internal/middleware"
	"github.com/Sofja96/gophermarket.git/internal/services"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"time"
)

type APIServer struct {
	echo           *echo.Echo
	address        string
	logger         zap.SugaredLogger
	accrualService *services.AccrualService
	store          *pg.Postgres
	OrdersChan     chan string
}

func New(ctx context.Context) *APIServer {
	var err error
	a := &APIServer{}
	c := config.LoadConfig()
	config.ParseFlags(c)
	helpers.NewLogger()

	a.address = c.Address
	a.echo = echo.New()
	//var wg sync.WaitGroup
	//orderchan := make(chan string)
	//var store storage.Storage
	a.store, err = pg.NewStorage(ctx, c.DatabaseDSN)
	if err != nil {
		log.Print(err)
	}
	log.Println(c.DatabaseDSN)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	//	key := c.HashKey
	a.accrualService = services.NewAccrualService(c.AccrualAddress, a.store)
	//	userHandler := handlers.NewUserHandler(store, a.accrualService)
	//userHandler := handlers.NewUserHandler(app.storage, app.accrualService, app.pool)
	a.logger = *logger.Sugar()
	a.echo.Use(middleware.WithLogging(a.logger))

	//a.accrualService.UpdateOrderStatus()

	//go a.accrualService.CheckOrderStatus()
	go a.accrualService.UpdateOrdersStatus(a.OrdersChan)
	//	a.accrualService.GetStatusOrder()

	//	a.checkOrderStatus(orderchan)

	group := a.echo.Group("api/user")
	{
		group.POST("/register", handlers.RegisterUser(a.store))
		group.POST("/login", handlers.LoginUser(a.store))
		group.Use(middleware.ValidateUser())
		{
			group.POST("/orders", handlers.PostOrder(a.store, a.OrdersChan))
			group.GET("/orders", handlers.GetOrders(a.store))
			group.GET("/balance", handlers.GetBalance(a.store))
			group.POST("/balance/withdraw", handlers.WithdrawBalance(a.store))
			group.GET("/withdrawals", handlers.GetWithdrawals(a.store))
			//	log.Println(<-orderchan)
			//	log.Println(a.OrdersChan)
			//userAPI.GET("/orders", userHandler.GetOrders)
			//userAPI.GET("/balance", userHandler.GetBalance)
			//userAPI.POST("/balance/withdraw", userHandler.WithdrawBalance)
			//userAPI.GET("/withdrawals", userHandler.GetWithdrawals)
		}
	}
	//a.echo.Use(middleware.GzipMiddleware())
	//a.echo.POST("/update/", UpdateJSON(store))
	//a.echo.POST("/updates/", UpdatesBatch(store))
	//a.echo.POST("/value/", ValueJSON(store))
	//a.echo.GET("/", GetAllMetrics(store))
	//a.echo.GET("/value/:typeM/:nameM", ValueMetric(store))
	//a.echo.POST("/update/:typeM/:nameM/:valueM", Webhook(store))
	//a.echo.GET("/ping", Ping(store))
	//go startTask(a.OrdersChan)
	//wg.Wait()
	return a
}

//func (a *APIServer) checkOrderStatus(ordersChan chan string) {
//	var wg sync.WaitGroup
//	outCh := make(chan models.OrderAccrual, 10)
//	wg.Add(1)
//	go func() {
//		ticker := time.NewTicker(1 * time.Second)
//		defer ticker.Stop()
//		//ctx := context.Background()
//		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
//		//defer cancel()
//
//		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//		//defer cancel ()
//		//err := uh.storage.UpdateOrderStatus(order.ID, models.PROCESSING)
//		//if err != nil {
//		//	log.Println("update order status: %v", err)
//		//}
//		for range ticker.C {
//			for order := range ordersChan {
//				resp, err := a.accrualService.GetStatusAccrual(order)
//				if err != nil {
//					log.Println("submit order: CalcOrderAccrual: %v", err)
//					return
//				}
//
//				if resp.Status != string(models.PROCESSED) {
//					time.AfterFunc(2*time.Second, func() {
//						ordersChan <- order
//					})
//					return
//				}
//				log.Println(order, "order number")
//				log.Println(resp, "resp")
//				outCh <- resp
//				//if resp.Status != string(models.PROCESSED) {
//				//	time.AfterFunc(5*time.Second, func() {
//				//		//ordersChan <- order
//				//	})
//				//	return
//				//}
//				//log.Println(resp.Order, "resp.order")
//				//if resp.Status == string(models.PROCESSED) {
//				//	log.Println(order, "processed")
//				//	//close(a.OrdersChan)
//				//
//				//}
//
//				//	close(outCh)
//
//				//close(a.OrdersChan)
//				//close(a.OrdersChan)
//				//	a.OrdersChan <- order
//				//err = uh.storage.UpdateOrderAccrualAndUserBalance(ctx, order.ID, userID, resp)
//				//if err != nil {
//				//	log.Println("submit order: UpdateOrderAccrualAndUserBalance: %v", err)
//				//	return
//				//}
//				//err = uh.storage.UpdateOrderStatus(ctx, order.ID, models.PROCESSED)
//				//if err != nil {
//				//	log.Println("update order status: %v", err)
//				//}
//				//	}
//
//			}
//		}
//		wg.Done()
//
//	}()
//	//wg.Add(1)
//	//go func() {
//	//	ticker := time.NewTicker(time.Duration(1) * time.Second)
//	//	defer ticker.Stop()
//	//	newOrders, err := a.store.GetOrderStatus(models.NEW)
//	//	if err != nil {
//	//		if err != nil {
//	//			log.Println("error get status order", err)
//	//
//	//		}
//	//		for range ticker.C {
//	//			//	var ordernum models.OrderAccrual
//	//			//var ordernum string
//	//			for _, ordernum := range newOrders {
//	//				outCh <- ordernum
//	//			}
//	//		}
//	//	}
//	//	wg.Done()
//	//}()
//	wg.Add(1)
//	go func() {
//		for order := range outCh {
//			log.Println(order.Order, "resp.order")
//			if order.Status == string(models.PROCESSED) || order.Status == string(models.INVALID) {
//				log.Println(order.Status, "order.status")
//				log.Println(order.Accrual, "accrual for begin update order")
//				err := a.store.UpdateOrder(order.Order, order.Status, order.Accrual)
//				if err != nil {
//					log.Println("error update order", err)
//					return
//				}
//				log.Println(order.Accrual, "accrual after update")
//				//close(a.OrdersChan)
//
//			}
//			if order.Status == string(models.PROCESSING) || order.Status == string(models.REGISTERED) {
//				log.Println(order.Status, "order.status")
//				err := a.store.UpdateOrderStatus(order.Order, models.OrderStatus(order.Status))
//				if err != nil {
//					log.Println("error update order", err)
//					return
//				}
//				//close(a.OrdersChan)
//
//			}
//		}
//		wg.Done()
//	}()
//	//	return outCh
//	//go startTask(ordersChan)
//	//wg.Wait()
//}

//	func (a *APIServer) UpdateStatusandBalace(ordersChan <-chan string) {
//		var wg sync.WaitGroup
//		outCh := make(chan models.OrderAccrual, 10)
//		wg.Add(1)
//		go func() {
//			ticker := time.NewTicker(1 * time.Second)
//			defer ticker.Stop()
//			//ctx := context.Background()
//			//ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
//			//defer cancel()
//
//			//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//			//defer cancel ()
//			//err := uh.storage.UpdateOrderStatus(order.ID, models.PROCESSING)
//			//if err != nil {
//			//	log.Println("update order status: %v", err)
//			//}
//			for range ticker.C {
//				for order := range ordersChan {
//					resp, err := a.accrualService.GetStatusAccrual(order)
//					if err != nil {
//						log.Println("submit order: CalcOrderAccrual: %v", err)
//						log.Println(order, "order in go func")
//						return
//					}
//
//					log.Println(order, "order number")
//					log.Println(resp, "resp")
//					//if resp.Status != string(models.PROCESSED) {
//					//	time.AfterFunc(5*time.Second, func() {
//					//		//ordersChan <- order
//					//	})
//					//	return
//					//}
//					log.Println(resp.Order, "resp.order")
//					if resp.Status == string(models.PROCESSED) {
//						log.Println(order, "processed")
//						//close(a.OrdersChan)
//
//					}
//					outCh <- resp
//					//close(a.OrdersChan)
//					//close(a.OrdersChan)
//					//	a.OrdersChan <- order
//					//err = uh.storage.UpdateOrderAccrualAndUserBalance(ctx, order.ID, userID, resp)
//					//if err != nil {
//					//	log.Println("submit order: UpdateOrderAccrualAndUserBalance: %v", err)
//					//	return
//					//}
//					//err = uh.storage.UpdateOrderStatus(ctx, order.ID, models.PROCESSED)
//					//if err != nil {
//					//	log.Println("update order status: %v", err)
//					//}
//					//	}
//
//				}
//			}
//			wg.Done()
//
//		}()
//		//	return outCh
//		//go startTask(ordersChan)
//		//wg.Wait()
//	}
func startTask(taskChan chan string) {
	for {
		select {
		case <-taskChan:
			return
		default:
			//	log.Println("Задача выполняется...")
		}
	}
}

func (a *APIServer) Start() error {
	go func() {
		err := a.echo.Start(a.address)
		if err != nil {
			log.Fatal(err)
		}
		helpers.Infof("starting the GopherMart server...", a.address)
		//log.Println("Running server on", a.address)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.echo.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	return nil
}
