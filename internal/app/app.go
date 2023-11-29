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
	a.store, err = pg.NewStorage(ctx, c.DatabaseDSN)
	if err != nil {
		helpers.Error("error creation storage: %s", err)
	}
	//TODO delete log
	helpers.Debug(c.DatabaseDSN)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	a.accrualService = services.NewAccrualService(c.AccrualAddress, a.store)
	a.logger = *logger.Sugar()
	a.echo.Use(middleware.WithLogging(a.logger))
	a.echo.Use(middleware.GzipMiddleware())

	go a.accrualService.UpdateOrdersStatus()

	group := a.echo.Group("api/user")
	{
		group.POST("/register", handlers.RegisterUser(a.store))
		group.POST("/login", handlers.LoginUser(a.store))
		group.Use(middleware.ValidateUser())
		{
			group.POST("/orders", handlers.PostOrder(a.store))
			group.GET("/orders", handlers.GetOrders(a.store))
			group.GET("/balance", handlers.GetBalance(a.store))
			group.POST("/balance/withdraw", handlers.WithdrawBalance(a.store))
			group.GET("/withdrawals", handlers.GetWithdrawals(a.store))
		}
	}
	return a
}

func (a *APIServer) Start() error {
	go func() {
		err := a.echo.Start(a.address)
		if err != nil {
			helpers.Fatal("error start GopherMart %s", err)
		}
		helpers.Infof("starting the GopherMart server...%s", a.address)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.echo.Shutdown(ctx); err != nil {
		helpers.Fatal("error shutting down gracefully %s", err)
	}
	return nil
}
