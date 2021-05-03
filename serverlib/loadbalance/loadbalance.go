package loadbalance

import (
	"kisgateway/serverlib/store/gateway"
	"sync"
)

var LbHandler *Lb

type LbItem struct {
	Name string
	Lb *LoadBalance
}

type Lb struct {
	LbMap map[string]*LbItem
	LbSlice []*LbItem
	lock sync.RWMutex
}

func NewLb() *Lb {
	return &Lb{
		LbMap: map[string]*LbItem{},
		LbSlice: []*LbItem{},
		lock:    sync.RWMutex{},
	}
}

func init()  {
	LbHandler = NewLb()
}

func (lb *Lb)GetLoadBalance(s *gateway.ServiceInfo) (*LoadBalance, error) {
	if v, ok := lb.LbMap[s.ServerInfo.ServiceName]; ok {
		return v.Lb, nil
	}

	return nil, nil
}