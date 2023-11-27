package main

import (
	"context"
	"github.com/Sofja96/gophermarket.git/internal/app"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	s := app.New(ctx)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	go func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.Shutdown(c)
	}()

	//defer s.Shutdown()
}
