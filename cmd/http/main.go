package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tech_check/internal/app"
	"tech_check/internal/handler"
	"time"
)

func main() {
	app := app.MustNew()
	server := setup(app)

	errChan := make(chan error, 1)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go start(server, app, errChan)
	wait(errChan, stopChan)
	shutdown(server)
}

func setup(app *app.App) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", app.Cfg.HTTP.Port),
		Handler: handler.New(app),
	}
}

func start(server *http.Server, app *app.App, errChan chan<- error) {
	log.Printf("starting app on port: %s\n", app.Cfg.HTTP.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		errChan <- err
	}
}

func wait(errChan <-chan error, stopChan <-chan os.Signal) {
	select {
	case <-stopChan:
		log.Println("stopping app...")
	case err := <-errChan:
		panic(err)
	}
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("server stopped gracefully")
}
