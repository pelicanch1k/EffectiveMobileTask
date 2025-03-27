package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/app"
)

// @title Song API
// @version 1.0
// @description API Server for EffectiveMobileTestTask

// @host localhost:80
func main() {
	app := app.NewApp()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		app.Stop(context.Background())
	}()

	app.Run()
}
