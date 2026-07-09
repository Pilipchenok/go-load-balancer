package balancer

import (
	"time"
	"context"
	"net"
	"log"
	"go-load-balancer/internal/backend"
)

func (bl *Balancer) RunHealthCheck(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		return
	}
	bl.checkAll()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			bl.checkAll()
		case <-ctx.Done():
			return
		}
	}
}

func (bl *Balancer) checkAll() {
	for i := 0; i < len(bl.backends); i++ {
		oldAlive := bl.backends[i].IsAlive()
		newAlive := bl.checkOne(bl.backends[i])
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
