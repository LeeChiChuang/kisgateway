package gateway

import "strings"

func GetIpList(s *ServiceInfo) []string {
	return strings.Split(s.LoadBalance.IpList, ";")
}

func GetWeightList(s *ServiceInfo) []string {
	return strings.Split(s.LoadBalance.WeightList, ";")
}

func GetBlackList(s *ServiceInfo) []string {
	return strings.Split(s.AccessControl.BlackList, ";")
}

func GetWhiteList(s *ServiceInfo) []string {
	return strings.Split(s.AccessControl.WhiteList, ";")
}