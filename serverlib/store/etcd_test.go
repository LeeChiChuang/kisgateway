package store

import (
	"fmt"
	"kisgateway/serverlib/store/gateway"
	"testing"
)

func TestEtcdStore_AddHttpService(t *testing.T) {
	e := NewEtcdStore("test")
	item := &gateway.ServiceInfo{
		ServerInfo:           &gateway.ServerInfo{
			Id:                   1,
			LoadType:             1,
			ServiceName:          "test",
			ServiceDesc:          "test server",
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
		AccessControl:        &gateway.AccessControl{
			OpenAuth:             1,
			WhiteList:            "127.0.0.8",
			BlackList:            "",
			WhiteHostName:        "",
			ClientipFlowLimit:    100,
			ServiceFlowLimit:     100,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
		LoadBalance:          &gateway.LoadBalance{
			CheckMethod:            0,
			CheckTimeout:           0,
			CheckInterval:          0,
			RoundType:              0,
			IpList:                 "",
			WeightList:             "",
			ForbidList:             "",
			UpstreamConnectTimeout: 0,
			UpstreamHeaderTimeout:  0,
			UpstreamIdleTimeout:    0,
			UpstreamMaxIdle:        0,
			XXX_NoUnkeyedLiteral:   struct{}{},
			XXX_unrecognized:       nil,
			XXX_sizecache:          0,
		},
		HttpRule: &gateway.HttpRule{
			RuleType:             0,
			Rule:                 "127.0.0.1",
			NeedHttps:            0,
			NeedStripUri:         0,
			NeedWebsocket:        0,
			UrlRewrite:           "",
			HeaderTransform:      "",
		},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	err := e.AddHttpService(item)
	if err != nil {
		t.Errorf("add http service error")
	}
}

func TestEtcdStore_GetHttpServices(t *testing.T) {
	e := NewEtcdStore("test")
	err := e.GetHttpServices()
	if err != nil {
		t.Errorf("get fail %+v \n", e)
		return
	}
	fmt.Printf("info list: %+v \n", e.HttpServices)
}