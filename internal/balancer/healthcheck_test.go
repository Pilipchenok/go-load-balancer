package balancer

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"bytes"
	"log"
	"strings"
	"context"
	"time"
	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/strategy"
)

func TestCheckOne(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer fakeServer.Close()

	fakeback, _ := backend.NewBackend(fakeServer.URL)
	lb := NewBalancer(
		[]*backend.Backend{fakeback},
		&strategy.RoundRobin{},
		&http.Client{},
	)
	res := lb.checkOne(fakeback)
	if res != true {
		t.Error("Working backend is dead")
	}
	if !fakeback.IsAlive() {
		t.Error("Backend must be marked as alive")
	}

	fakeLostServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"status": "bad gateway"}`))
	}))
	fakeLostServer.Close()

	fakeLostBack, _ := backend.NewBackend(fakeLostServer.URL)
	lb2 := NewBalancer(
		[]*backend.Backend{fakeLostBack},
		&strategy.RoundRobin{},
		&http.Client{},
	)
	res2 := lb2.checkOne(fakeLostBack)
	if res2 != false {
		t.Error("Not working backend works")
	}
	if fakeLostBack.IsAlive() {
		t.Error("Bad backend must be marked as dead")
	}
}

func TestCheckAll(t *testing.T) {
	serverAlive := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer serverAlive.Close()

	serverDead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	serverDead.Close() 

	back1, _ := backend.NewBackend(serverAlive.URL)
	back2, _ := backend.NewBackend(serverDead.URL)
	
	back1.SetAlive(false)
	back2.SetAlive(true)

	backends := []*backend.Backend{
		back1,
		back2,
	}

	bl := &Balancer{backends: backends}

	var logBuf bytes.Buffer
	oldWriter := log.Writer()
	log.SetOutput(&logBuf)
	defer log.SetOutput(oldWriter)

	bl.checkAll()

	output := logBuf.String()
	
	if !strings.Contains(output, "backend0: false -> true") {
		t.Errorf("Ожидался лог для backend0, получено:\n%s", output)
	}
	if !strings.Contains(output, "backend1: true -> false") {
		t.Errorf("Ожидался лог для backend1, получено:\n%s", output)
	}
}

func TestRunHealthCheck(t *testing.T) {
	serverAlive := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer serverAlive.Close()

	serverDead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	serverDead.Close() 

	back1, _ := backend.NewBackend(serverAlive.URL)
	back2, _ := backend.NewBackend(serverDead.URL)
	
	back1.SetAlive(false)
	back2.SetAlive(true)

	backends := []*backend.Backend{
		back1,
		back2,
	}

	bl := NewBalancer(backends, &strategy.RoundRobin{}, &http.Client{})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	done := make(chan struct{})
	go func() {
		bl.RunHealthCheck(ctx, time.Second * 1)
		close(done)
	}()
	cancel()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		t.Fatal("RunHealthCheck did not stop")
	}
}
