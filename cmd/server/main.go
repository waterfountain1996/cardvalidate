package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/waterfountain1996/cardvalidate/api"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("POST /validate", api.ValidationHandler())

	httpSrv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	go func() {
		log.Printf("Serving on %s\n", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown(): %s\n", err)
	}
}
