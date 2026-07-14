package main

import (
	"log"
	"net/http"
	"time"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/balancer"
	"go-load-balancer/internal/config"
	"go-load-balancer/internal/strategy"
)

func main() {
	configPath := "configs/config.yaml"
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Import config error: %v", err)
	}

	confBacks := cfg.Backends
	backends := []*backend.Backend{}
	for i := 0; i < len(confBacks); i++ {
		back, err := backend.NewBackend(confBacks[i].URL)
		if err != nil {
			log.Fatalf("Backend create error: %v", err)
		}
		backends = append(backends, back)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var myBalancer *balancer.Balancer
	if cfg.Strategy == "round-robin" {
		myBalancer = balancer.NewBalancer(
			backends,
			&strategy.RoundRobin{},
			client,
		)
	} else if cfg.Strategy == "least-conn" {
		myBalancer = balancer.NewBalancer(
			backends,
			&strategy.LeastConn{},
			client,
		)
	} else {
		log.Fatal("Unknown strategy")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		myBalancer.RunHealthCheck(ctx, cfg.HealthCheckInterval)
	}()

	server := &http.Server{
    Addr: fmt.Sprintf(":%d", cfg.Port),
    Handler: myBalancer,
	}

	go func() {
		log.Printf("Balancer started on :%d", cfg.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	errSD := server.Shutdown(shutdownCtx)
	if errSD != nil {
		log.Fatalf("Stop server error: %v", errSD)
	}

	log.Println("Server stopped")
}
