package strategy

import (
	"testing"
	"go-load-balancer/internal/backend"
)

func TestRoundRobin(t *testing.T) {
	back1, _ := backend.NewBackend("http://localhost:8080")
	back2, _ := backend.NewBackend("http://localhost:8081")
	back3, _ := backend.NewBackend("http://localhost:8082")
	backs := []*backend.Backend{back1, back2, back3}
	rr := RoundRobin{}
	chooseBack := rr.Next(backs)
	if chooseBack != back1 {
		t.Error("Not needed backend")
	}
	chooseBack = rr.Next(backs)
	if chooseBack != back2 {
		t.Error("Not needed backend")
	}
	chooseBack = rr.Next(backs)
	if chooseBack != back3 {
		t.Error("Not needed backend")
	}
	chooseBack = rr.Next(backs)
	if chooseBack != back1 {
		t.Error("Not needed backend")
	}
}

func TestDeadRoundRobin(t *testing.T) {
	back1, _ := backend.NewBackend("http://localhost:8080")
	back2, _ := backend.NewBackend("http://localhost:8081")
	back3, _ := backend.NewBackend("http://localhost:8082")
	back1.SetAlive(false)
	backs := []*backend.Backend{back1, back2, back3}
	rr := RoundRobin{}
	chooseBack := rr.Next(backs)
	if chooseBack != back2 {
		t.Error("Not needed backend")
	}
	back2.SetAlive(false)
	back3.SetAlive(false)
	chooseBack = rr.Next(backs)
	if chooseBack != nil {
		t.Error("Return dead backend")
	}
}

func TestEmptyRoundRobin(t *testing.T) {
	backs := []*backend.Backend{}
	rr := RoundRobin{}
	chooseBack := rr.Next(backs)
	if chooseBack != nil {
		t.Error("Return with no backends")
	}
}

func TestLeastConn(t *testing.T) {
	back1, _ := backend.NewBackend("http://localhost:8080")
	back2, _ := backend.NewBackend("http://localhost:8081")
	back3, _ := backend.NewBackend("http://localhost:8082")
	back1.IncrementConnections()
	back1.IncrementConnections()
	back3.IncrementConnections()
	backs := []*backend.Backend{back1, back2, back3}
	lc := LeastConn{}
	chooseBack := lc.Next(backs)
	if chooseBack != back2 {
		t.Error("Not needed backend")
	}
	back2.IncrementConnections()
	back2.IncrementConnections()
	chooseBack = lc.Next(backs)
	if chooseBack != back3 {
		t.Error("Not needed backend")
	}
	back3.IncrementConnections()
	back3.IncrementConnections()
	back2.IncrementConnections()
	chooseBack = lc.Next(backs)
	if chooseBack != back1 {
		t.Error("Not needed backend")
	}
}

func TestDeadLeastConn(t *testing.T) {
	back1, _ := backend.NewBackend("http://localhost:8080")
	back2, _ := backend.NewBackend("http://localhost:8081")
	back3, _ := backend.NewBackend("http://localhost:8082")
	back1.SetAlive(false)
	back2.IncrementConnections()
	backs := []*backend.Backend{back1, back2, back3}
	lc := LeastConn{}
	chooseBack := lc.Next(backs)
	if chooseBack != back3 {
		t.Error("Not needed backend")
	}
	back2.SetAlive(false)
	back3.SetAlive(false)
	chooseBack = lc.Next(backs)
	if chooseBack != nil {
		t.Error("Return dead backend")
	}
}

func TestEmptyLeastConn(t *testing.T) {
	backs := []*backend.Backend{}
	lc := LeastConn{}
	chooseBack := lc.Next(backs)
	if chooseBack != nil {
		t.Error("Return with no backends")
	}
}
