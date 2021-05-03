package loadbalance

type RoundRobinBalance struct {
	
}

func (r RoundRobinBalance) Add(...string) error {
	panic("implement me")
}

func (r RoundRobinBalance) Get() (string, error) {
	panic("implement me")
}
 
