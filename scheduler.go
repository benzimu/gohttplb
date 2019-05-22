package gohttplb

import (
	"gohttplb/utils"
	"log"
	"net/http"
)

// SchedPolicy a schedule policy for Request
type SchedPolicy int

const (
	// PolicyRandom random policy
	PolicyRandom SchedPolicy = iota + 1
	// PolicyOrder order policy
	PolicyOrder
	// PolicyFlow flow policy
	PolicyFlow
	// PolicyFlowWeight flow weight policy
	PolicyFlowWeight
)

// NewRScheduler new RScheduler interface
func NewRScheduler(r *R) RScheduler {
	switch r.SchedPolicy {
	case PolicyRandom:
		return &RandomScheduler{r}
	case PolicyOrder:
		return &OrderScheduler{r}
	default:
		return &RandomScheduler{r}
	}
}

// RScheduler a scheduler for Request
type RScheduler interface {
	schedule(*rArgs) (*http.Response, error)
}

// RandomScheduler for PolicyRandom
type RandomScheduler struct {
	*R
}

func (sched *RandomScheduler) schedule(rargs *rArgs) (resp *http.Response, err error) {
	for retry := 0; retry < sched.Retry; retry++ {
		servers := sched.servers
		for count := 0; count < len(servers); count++ {
			index := utils.GenRandIntn(len(servers))
			server := servers[index]
			rargs.url = server + rargs.url
			resp, err = sched.do(rargs)
			if err != nil {
				log.Printf("WARN: %s; server: %s", err, server)
				servers = append(servers[:index], servers[index+1:]...)
				continue
			}
			return
		}
	}
	return
}

// OrderScheduler for PolicyOrder scheduler
type OrderScheduler struct {
	*R
}

func (sched *OrderScheduler) schedule(rargs *rArgs) (resp *http.Response, err error) {
	for retry := 0; retry < sched.Retry; retry++ {
		for _, server := range sched.servers {
			rargs.url = server + rargs.url
			resp, err = sched.do(rargs)
			if err != nil {
				log.Printf("WARN: %s; server: %s", err, server)
				continue
			}
			return
		}
	}
	return
}

// FlowScheduler for PolicyFlow scheduler
type FlowScheduler struct {
	*R
}

func (sched *FlowScheduler) schedule(rargs *rArgs) (resp *http.Response, err error) {
	for retry := 0; retry < sched.Retry; retry++ {
		for _, server := range sched.servers {
			rargs.url = server + rargs.url
			resp, err = sched.do(rargs)
			if err != nil {
				log.Printf("WARN: %s; server: %s", err, server)
				continue
			}
			return
		}
	}
	return
}
