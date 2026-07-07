package strategy

import (
	"go-load-balancer/internal/backend"
)

type LeastConn struct {

}

func (lc *LeastConn) Next(backends []*backend.Backend) *backend.Backend {
	if len(backends) == 0 {
		return nil
	}
	minConn := int64(0)
	var needBack *backend.Backend
	for i := 0; i < len(backends); i++ {
		if needBack == nil && backends[i].IsAlive() {
			needBack = backends[i]
			minConn = backends[i].ActiveConnections()
		}
		if backends[i].ActiveConnections() < minConn && backends[i].IsAlive() {
			minConn = backends[i].ActiveConnections()
			needBack = backends[i]
		}
	}
	return needBack
}
