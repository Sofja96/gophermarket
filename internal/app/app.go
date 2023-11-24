package app

import (
	"context"
	"github.com/Sofja96/gophermarket.git/internal/app/config"
	"github.com/Sofja96/gophermarket.git/internal/handlers"
	"github.com/Sofja96/gophermarket.git/internal/middleware"
	"github.com/Sofja96/gophermarket.git/internal/storage/pg"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"log"
)

type APIServer struct {
	echo    *echo.Echo
	address string
	logger  zap.SugaredLogger
}

func New(ctx context.Context) *APIServer {
	a := &APIServer{}
	c := config.LoadConfig()
	config.ParseFlags(c)

	a.address = c.Address
	a.echo = echo.New()
	//var store storage.Storage
	var err error
	store, err := pg.NewStorage(ctx, c.DatabaseDSN)
	if err != nil {
		log.Print(err)
	}
	log.Println(c.DatabaseDSN)
	//if len(c.DatabaseDSN) == 0 {
	//	store, err = memory.NewInMemStorage(c.StoreInterval, c.FilePath, c.Restore)
	//	if err != nil {
	//		log.Print(err)
	//	}
	//} else {
	//	store, err = database.NewStorage(c.DatabaseDSN)
	//}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	//	key := c.HashKey
	//as := services.NewAccrualService(c.AccrualAddress)
	//userHandler := handlers.NewUserHandler(store, as)
	//userHandler := handlers.NewUserHandler(app.storage, app.accrualService, app.pool)
	a.logger = *logger.Sugar()
	a.echo.Use(middleware.WithLogging(a.logger))
	//if len(key) != 0 {
	//	a.echo.Use(middleware.HashMacMiddleware([]byte(key)))
	//}
	group := a.echo.Group("api/user")
	{
		group.POST("/register", handlers.RegisterUser(store))
		group.POST("/login", handlers.LoginUser(store))
		group.Use(middleware.ValidateUser())

	}
	//{
	//	group.POST("/orders", handlers.PostOrder(userHandler))
	//	//userAPI.GET("/orders", userHandler.GetOrders)
	//	//userAPI.GET("/balance", userHandler.GetBalance)
	//	//userAPI.POST("/balance/withdraw", userHandler.WithdrawBalance)
	//	//userAPI.GET("/withdrawals", userHandler.GetWithdrawals)
	//}
	//a.echo.Use(middleware.GzipMiddleware())
	//a.echo.POST("/update/", UpdateJSON(store))
	//a.echo.POST("/updates/", UpdatesBatch(store))
	//a.echo.POST("/value/", ValueJSON(store))
	//a.echo.GET("/", GetAllMetrics(store))
	//a.echo.GET("/value/:typeM/:nameM", ValueMetric(store))
	//a.echo.POST("/update/:typeM/:nameM/:valueM", Webhook(store))
	//a.echo.GET("/ping", Ping(store))
	return a
}

func (a *APIServer) Start() error {
	err := a.echo.Start(a.address)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Running server on", a.address)

	return nil
}
