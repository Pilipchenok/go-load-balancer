package strategy

import (
	"go-load-balancer/internal/backend"
	"sync/atomic"
)

type RoundRobin struct {
	n atomic.Uint32
}

func (rr *RoundRobin) Next(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}
	n := rr.n.Add(1) - 1
	l := uint32(len(backends))
	for i := uint32(0); i < l; i++ {
		k := (n + i) % l
		if backends[k].IsAlive() {
			return backends[k]
		}
	}
	return nil
}
