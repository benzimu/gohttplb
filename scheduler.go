package gohttplb

import (
	"log"
	"net/http"
	"sync"
)

// SchedPolicy a schedule policy for Request
type SchedPolicy int

const (
	// PolicyRandom random policy
	PolicyRandom SchedPolicy = iota + 1
	// PolicyOrder order policy
	PolicyOrder
	// PolicyFlowAverage flow average policy
	PolicyFlowAverage
)

var flowAverageScheduler *FlowAverageScheduler

// NewRScheduler new RScheduler interface
func NewRScheduler(r *R) RScheduler {
	switch r.SchedPolicy {
	case PolicyRandom:
		return &RandomScheduler{r}
	case PolicyOrder:
		return &OrderScheduler{r}
	case PolicyFlowAverage:
		if flowAverageScheduler == nil {
			flowAverageScheduler = &FlowAverageScheduler{R: r}
		}
		return flowAverageScheduler
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
			tmp := rargs
			index := GenRandIntn(len(servers))
			server := servers[index]
			tmp.url = server + tmp.url
			resp, err = sched.do(tmp)
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
			tmp := rargs
			tmp.url = server + tmp.url
			resp, err = sched.do(tmp)
			if err != nil {
				log.Printf("WARN: %s; server: %s", err, server)
				continue
			}
			return
		}
	}
	return
}

// FlowAverageScheduler for PolicyFlowAverages scheduler
type FlowAverageScheduler struct {
	*R
	fm sync.Map
}

func (sched *FlowAverageScheduler) clearFm() {
	for _, server := range sched.servers {
		sched.fm.Delete(server)
	}
}

func (sched *FlowAverageScheduler) getServer(errServers []string) (string, bool) {
	var hasV bool
	var minCount int
	index := -1
	// Select the service with the least flow count
	for i, server := range sched.servers {
		// except error servers
		if len(errServers) > 0 {
			if ExistStringSlice(server, errServers) {
				continue
			}
		}

		count, ok := sched.fm.Load(server)
		if !hasV && !ok {
			continue
		}

		if hasV && !ok {
			count = 0
		}

		countI := count.(int)
		if !hasV {
			hasV = true
			minCount = countI
			index = i
			continue
		}
		if countI < minCount {
			if countI+100 < minCount {
				go sched.clearFm()
			}
			minCount = countI
			index = i
		}
	}

	// Iterative once already
	if index == -1 {
		return sched.servers[0], true
	}

	if minCount > 2<<20 {
		go sched.clearFm()
	}

	return sched.servers[index], false
}

func (sched *FlowAverageScheduler) schedule(rargs *rArgs) (resp *http.Response, err error) {
	errServers := make([]string, 0)
	for retry := 0; retry < sched.Retry*len(sched.servers); retry++ {
		tmp := rargs
		server, once := sched.getServer(errServers)
		if once {
			errServers = errServers[:0]
		}
		tmp.url = server + tmp.url
		resp, err = sched.do(tmp)
		if err != nil {
			log.Printf("WARN: %s; server: %s", err, server)
			errServers = append(errServers, server)
			continue
		}
		// flow count increase 1 if the request success
		val, loaded := sched.fm.LoadOrStore(server, 1)
		if loaded {
			valI := val.(int)
			valI++
			sched.fm.Store(server, valI)
		}
		return
	}
	return
}
