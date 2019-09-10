package strategy

import (
	"gohttplb/utils"
)

// ServerItem is server item with weighted
type ServerItem struct {
	Server   string
	Weighted int
}

// WeightedRoundRobinMaker is Weighted Round Robin Balancing Algorithm for StrategyWeightedRoundRobin
type WeightedRoundRobinMaker struct {
	servers []ServerItem
	// len of servers
	n int
	// last server index
	index int
	// current weighted
	curW int
	// greatest common divisor of all weighted
	gcdW int
	// max weighted
	maxW int
}

// NewWeightedRoundRobinMaker new WeightedRoundRobinMaker and init
func NewWeightedRoundRobinMaker(servers []ServerItem) *WeightedRoundRobinMaker {
	if len(servers) == 0 {
		return nil
	}

	maker := new(WeightedRoundRobinMaker)
	for _, server := range maker.servers {
		if maker.gcdW == 0 {
			maker.n = len(maker.servers)
			maker.index = -1
			maker.curW = 0
			maker.gcdW = server.Weighted
			maker.maxW = server.Weighted
		} else {
			maker.gcdW = utils.Gcd(maker.gcdW, server.Weighted)
			if maker.maxW < server.Weighted {
				maker.maxW = server.Weighted
			}
		}
	}
	return maker
}

func (maker *WeightedRoundRobinMaker) Make() (server string) {
	if maker.n == 0 {
		return ""
	}

	if maker.n == 1 {
		return maker.servers[0].Server
	}

	for {
		maker.index = (maker.index + 1) % maker.n
		if maker.index == 0 {
			maker.curW = maker.curW - maker.gcdW
			if maker.curW <= 0 {
				maker.curW = maker.maxW
				if maker.curW == 0 {
					return ""
				}
			}
		}
		if maker.servers[maker.index].Weighted >= maker.curW {
			return maker.servers[maker.index].Server
		}
	}
}
