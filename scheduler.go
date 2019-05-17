package gohttplb

import (
	"gohttplb/utils"
	"net/http"
)

// SchedPolicy a schedule policy for Request
type SchedPolicy int

const (
	// PolicyRandom for random policy
	PolicyRandom SchedPolicy = iota + 1
	// PolicyRandomRetry for random retry policy
	PolicyRandomRetry
	// PolicyOrder for servers order policy
	PolicyOrder
	// PolicyOrderRetry for servers order retry policy
	PolicyOrderRetry
	// PolicyFlow for flow policy
	PolicyFlow
	// PolicyFlowWeight for flow weight policy
	PolicyFlowWeight
)

// NewRScheduler new RScheduler interface
func NewRScheduler(r *R) RScheduler {
	switch r.SchedPolicy {
	case PolicyRandom:
		return &RandomScheduler{r}
	case PolicyRandomRetry:
		return &RandomRetryScheduler{r}
	default:
		return &RandomRetryScheduler{r}
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

func (sche *RandomScheduler) schedule(rargs *rArgs) (resp *http.Response, err error) {
	server := sche.servers[utils.GenRandIntn(len(sche.servers))]
	rargs.url = server + rargs.url
	return sche.do(rargs)
}

// RandomRetryScheduler for PolicyRandomRetry
type RandomRetryScheduler struct {
	*R
}

func (sche *RandomRetryScheduler) schedule(rargs *rArgs) (resp *http.Response, err error) {
	for count := 1; count <= sche.Retry; count++ {
		server := sche.servers[utils.GenRandIntn(len(sche.servers))]
		rargs.url = server + rargs.url
	}

	return sche.do(rargs)
}
