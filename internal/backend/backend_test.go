package backend

import (
	"testing"
	"sync"
)

func TestNewBackend_ValidURL(t *testing.T) {
	back, err := NewBackend("http://localhost:8080")
	if err != nil {
		t.Error("Create backend error", err)
	}
	if back == nil {
		t.Error("Backend is nil")
	}
	if back.IsAlive() != true {
		t.Error("Backend is not alive")
	}
	if back.ActiveConnections() != 0 {
		t.Error("Backend active connections is not 0");
	}
}

func TestNewBackend_InvalidURL(t *testing.T) {
	invalidURLs := []string{
		"",
		"localhost:8080",
		"http://",
	}

	for _, badURL := range invalidURLs {
		t.Run(badURL, func(t *testing.T) {
			back, err := NewBackend(badURL)

			if err == nil || back != nil {
				t.Errorf("Error expected")
			}
		})
	}
}


func TestBackend_SetAlive(t *testing.T) {
	back, _ := NewBackend("http://localhost:8080")
	back.SetAlive(false)
	if back.IsAlive() != false {
		t.Error("SetAlive Error")
	}
	back.SetAlive(true)
	if back.IsAlive() != true {
		t.Error("SetAlive Error")
	}
}

func TestBackend_Connections(t *testing.T) {
	back, _ := NewBackend("http://localhost:8080")
	back.IncrementConnections()
	back.IncrementConnections()
	back.IncrementConnections()
	if back.ActiveConnections() != 3 {
		t.Error("IncrementConnections Error")
	}
	back.DecrementConnections()
	if back.ActiveConnections() != 2 {
		t.Error("DecrementConnections Error")
	}
}

func TestBackend_ConcurrentAccess(t *testing.T) {
	back, _ := NewBackend("http://localhost:8080")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			back.IncrementConnections()
			back.IsAlive()
			back.DecrementConnections()
		}()
	}
	wg.Wait()
	if back.activeConnections != 0 {
		t.Error("ConcurrentAccess Error")
	}
}
