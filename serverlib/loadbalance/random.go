package loadbalance

import (
	"errors"
	"math/rand"
)

var GetRandomLbFail = errors.New("get random lb fail")
type RandomLb struct {
	ss []string
}

func (lb *RandomLb)Add(ss ...string) error {
	lb.ss = append(lb.ss, ss...)
	return nil
}

func (lb *RandomLb)Get() (string, error) {
	if len(lb.ss) == 0 {
		return "", GetRandomLbFail
	}

	return lb.ss[rand.Intn(len(lb.ss))], nil
}
