package gohttplb

import "sync"

// RoundRobinMaker is Round Robin Balancing Algorithm for StrategyRoundRobin
type RoundRobinMaker struct {
	servers []string
	mutex   sync.Mutex
	next    int
}

// Make implement Scheduler interface
func (maker *RoundRobinMaker) Make() (server string) {
	// maker.mutex.Lock()
	// defer maker.mutex.Unlock()
	server = maker.servers[maker.next]
	maker.next = (maker.next + 1) % len(maker.servers)
	return
}
