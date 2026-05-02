package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
	"encoding/json"
	"strconv"
)

type Backend struct {
	url *url.URL
	mu *sync.RWMutex
	alive bool
	activeConnections int64
}

func NewBackend(newUrl string) (*Backend, error){
	parcedUrl, err := url.Parse(newUrl)
	if err != nil {
		return &Backend{}, err
	}
	return &Backend {
		parcedUrl,
		new(sync.RWMutex),
		true,
		0,
	}, nil
}

func (b *Backend) IsAlive () bool {
	b.mu.RLock()
	ans := b.alive
	b.mu.RUnlock()
	return ans
}

func (b *Backend) SetAlive (alive bool) {
	b.mu.Lock()
	b.alive = alive
	b.mu.Unlock()
}

func (b *Backend) URL () *url.URL{
	return b.url
}

func (b *Backend) GetActiveConnections() int64{
	ans := atomic.LoadInt64(&b.activeConnections)
	return ans
}

func (b *Backend) IncrementConnections() {
	atomic.AddInt64(&b.activeConnections, 1)
}

func (b *Backend) DecrementConnections() {
	atomic.AddInt64(&b.activeConnections, -1)
}

func (b *Backend) String() string {
	ans := map[string]string {
		"url": b.url.String(),
		"alive": strconv.FormatBool(b.alive),
		"active connections": strconv.Itoa(int(b.activeConnections)),
	}
	ansString, _ := json.Marshal(ans)
	return string(ansString)
}
