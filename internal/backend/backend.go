package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
	"fmt"
)

type Backend struct {
	url *url.URL
	mu sync.RWMutex
	alive bool
	activeConnections int64
}

func NewBackend(newURL string) (*Backend, error){
	parsedURL, err := url.Parse(newURL)
	if err != nil {
		return nil, fmt.Errorf("Parse URL Error %q: %w", newURL, err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid URL %q: scheme and host are required", newURL)
	}
	return &Backend {
		url: parsedURL,
		mu: sync.RWMutex{},
		alive: true,
		activeConnections: 0,
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

func (b *Backend) ActiveConnections() int64{
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
	return b.url.String()
}
