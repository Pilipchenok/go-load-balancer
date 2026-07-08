package balancer

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"slices"
	"go-load-balancer/internal/backend"
	"go-load-balancer/internal/strategy"
)

func TestNewBalancer(t *testing.T) {
	back1, _ := backend.NewBackend("http://localhost:8080")
	back2, _ := backend.NewBackend("http://localhost:8081")
	back3, _ := backend.NewBackend("http://localhost:8082")
	backs := []*backend.Backend{back1, back2, back3}
	st := strategy.RoundRobin{}
	cl := &http.Client {
		Timeout: 30,
	}
	bl := NewBalancer(backs, &st, cl)
	if !slices.Equal(bl.backends, backs) || bl.strategy != &st || bl.client != cl {
		t.Error("NewBalancer change data")
	}
}

func mustBackend(t *testing.T, rawURL string) *backend.Backend {
	t.Helper()
	b, err := backend.NewBackend(rawURL)
	if err != nil {
		t.Fatalf("Failed to create backend %q: %v", rawURL, err)
	}
	return b
}

func TestServeHTTP_ProxyToBackend(t *testing.T) {
	fakeBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer fakeBackend.Close()

	back := mustBackend(t, fakeBackend.URL)
	lb := NewBalancer(
		[]*backend.Backend{back},
		&strategy.RoundRobin{},
		&http.Client{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	lb.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	expected := `{"status": "ok"}`
	if rec.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %q", rec.Header().Get("Content-Type"))
	}
	if back.ActiveConnections() != 0 {
		t.Errorf("Expected 0 active connections, got %d", back.ActiveConnections())
	}
}

func TestServeHTTP_AllDeadReturns503(t *testing.T) {
	back1 := mustBackend(t, "http://localhost:8081")
	back2 := mustBackend(t, "http://localhost:8082")
	back1.SetAlive(false)
	back2.SetAlive(false)

	lb := NewBalancer(
		[]*backend.Backend{back1, back2},
		&strategy.RoundRobin{},
		&http.Client{},
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	lb.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}
}

func TestServeHTTP_BackendError(t *testing.T) {
	fakeBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer fakeBackend.Close()

	back := mustBackend(t, fakeBackend.URL)
	lb := NewBalancer(
		[]*backend.Backend{back},
		&strategy.RoundRobin{},
		&http.Client{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	rec := httptest.NewRecorder()
	lb.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
	if rec.Body.String() != "internal error" {
		t.Errorf("expected body 'internal error', got %q", rec.Body.String())
	}
}
