package main

import (
	"context"
	"github.com/Sofja96/gophermarket.git/internal/app"
	"github.com/Sofja96/gophermarket.git/internal/helpers"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	s := app.New(ctx)
	if err := s.Start(); err != nil {
		helpers.Fatal("error start service GopherMart %s", err)
	}
}
