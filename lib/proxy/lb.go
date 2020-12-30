package proxy

import (
	"errors"
	"sync"
)

var (
	// ErrEmptyTargets is used when no targets are provided
	ErrEmptyTargets = errors.New("no targets are provided")

	// ErrUnsupportedStrategy is used when an unsupported LB strategy is provided
	ErrUnsupportedStrategy = errors.New("load balancing strategy not supported")
)

// LB holds the methods which the load balancing algorithms must implement
type LB interface {
	Elect(hosts []string) (string, error)
}

// NewLB creates a balancer as per the
func NewLB(strategy string) (LB, error) {
	if strategy != "rr" {
		return nil, ErrUnsupportedStrategy
	}
	return newRoundrobinLB(), nil
}

type roundrobinLB struct {
	current int
	mu      sync.RWMutex
}

func newRoundrobinLB() *roundrobinLB {
	return &roundrobinLB{}
}

func (r *roundrobinLB) Elect(hosts []string) (string, error) {
	if len(hosts) == 0 {
		return "", ErrEmptyTargets
	}

	if len(hosts) == 1 {
		return hosts[0], nil
	}

	// reset the current host position once all the hosts are used
	if r.current >= len(hosts) {
		r.current = 0
	}

	host := hosts[r.current]
	r.mu.Lock()
	r.current++
	r.mu.Unlock()

	return host, nil
}
