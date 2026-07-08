package balancer

import (
	"net/http"
	"io"
	"net"
	"log"
	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/strategy"
)

type Balancer struct {
	backends []*backend.Backend
	strategy strategy.Strategy
	client *http.Client
}

func New(backs []*backend.Backend, st strategy.Strategy, cl *http.Client) (*Balancer) {
	return &Balancer {
		backends: backs,
		strategy: st,
		client: cl,
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	back := b.strategy.Next(b.backends)
	if back == nil {
		http.Error(w, "503 Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	back.IncrementConnections()
	defer back.DecrementConnections()
	backendURL := back.URL().Scheme + "://" + back.URL().Host + r.URL.Path
	if r.URL.RawQuery != "" {
		backendURL += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequestWithContext(
		r.Context(),
		r.Method,
		backendURL,
		r.Body,
	)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	for key, values := range r.Header {
		for _, newV := range values {
			req.Header.Add(key, newV)
		}
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIP = r.RemoteAddr 
	}
	req.Header.Add("X-Real-IP", clientIP)
	req.Header.Add("X-Forwarded-For", clientIP)
	resp, err := b.client.Do(req)
	if err != nil {
		http.Error(w, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, newV := range values {
			w.Header().Add(key, newV)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, errCopy := io.Copy(w, resp.Body)
	if errCopy != nil {
		log.Printf("Copy answer error: %v", errCopy)
		return
	}
}
