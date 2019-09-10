package strategy

// LoadBalancingStrategy is load balancing strategy for selecting request server
type LoadBalancingStrategy int

const (
	// StrategyRoundRobin is Round Robin Scheduling strategy
	StrategyRoundRobin LoadBalancingStrategy = iota + 1
	// StrategyWeightedRoundRobin is Weighted Round Robin Scheduling strategy
	StrategyWeightedRoundRobin
)

// Scheduler make a valid server for Request
type Scheduler interface {
	Make() string
}

// NewScheduler new strategy Scheduler
func NewScheduler(strategy LoadBalancingStrategy, servers []string, serverWeighteds []ServerItem) Scheduler {
	if len(servers) == 0 && len(serverWeighteds) == 0 {
		return nil
	}

	var scheduler Scheduler
	switch strategy {
	case StrategyRoundRobin:
		scheduler = &RoundRobinMaker{servers: servers}
	case StrategyWeightedRoundRobin:
		scheduler = NewWeightedRoundRobinMaker(serverWeighteds)
	default:
		scheduler = &RoundRobinMaker{servers: servers}
	}
	return scheduler
}
