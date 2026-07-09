package balancer

import (
	"time"
	"context"
	"net"
	"log"
	"go-load-balancer/internal/backend"
)

func (bl *Balancer) RunHealthCheck(backends []*backend.Backend, ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			bl.checkAll(backends)
		case <-ctx.Done():
			return
		}
	}
}

func (bl *Balancer) checkAll(backends []*backend.Backend) {
	for i := 0; i < len(backends); i++ {
		oldAlive := backends[i].IsAlive()
		newAlive := bl.checkOne(backends[i])
		if oldAlive != newAlive {
			log.Printf("backend%d: %v -> %v", i, oldAlive, newAlive)
		}
	}
}

func (bl *Balancer) checkOne(backend *backend.Backend) bool{
	conn, err := net.DialTimeout("tcp", backend.URL().Host, time.Second)
	if err != nil {
		backend.SetAlive(false)
		return false
	} else {
		backend.SetAlive(true)
		conn.Close()
		return true
	}
}
