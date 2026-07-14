package main

import (
	"log"
	"net/http"
	"time"
	"context"

	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/balancer"
	"go-load-balancer/internal/config"
	"go-load-balancer/internal/strategy"
)

func main() {
	configPath := "configs/config.yaml"
	config, err := config.Load(configPath)
	if err != nil {
		log.Printf("Import config error: %v", err)
		return
	}

	backends := []*backend.Backend{}
	confBacks := config.Backends
	for i := 0; i < len(confBacks); i++ {
		back, err := backend.NewBackend(confBacks[i].URL)
		if err != nil {
			log.Printf("Backend create error: %v", err)
			return
		}
		backends[i] = back
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var myBalancer balancer.Balancer
	if config.Strategy == "round_robin" {
		myBalancer = *balancer.NewBalancer(
			backends,
			&strategy.RoundRobin{},
			&client,
		)
	} else if config.Strategy == "least_conn" {
		myBalancer = *balancer.NewBalancer(
			backends,
			&strategy.LeastConn{},
			&client,
		)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		myBalancer.RunHealthCheck(ctx, time.Second * 1)
		cancel()
	}()

	var w http.ResponseWriter
	var r *http.Request
	myBalancer.ServeHTTP(w, r)
}
