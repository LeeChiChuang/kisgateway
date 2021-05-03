package loadbalance

type Hash func(data []byte) uint32

type ConsistentHashBalance struct {

}

func NewConsistentHashBalance(n int, fn Hash) *ConsistentHashBalance {
	return &ConsistentHashBalance{}
}

func (r ConsistentHashBalance) Add(...string) error {
	panic("implement me")
}

func (r ConsistentHashBalance) Get() (string, error) {
	panic("implement me")
}

