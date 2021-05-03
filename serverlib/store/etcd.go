package store

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"kisgateway/serverlib/logx"
	"kisgateway/serverlib/store/gateway"
	"time"
)

type EtcdStore struct {
	prefix     string
	HttpPrefix string
	HttpServices []*gateway.HttpService
	rawClient  *clientv3.Client
}

type Services struct {
	HttpServicesM map[string]*gateway.HttpService
	HttpService *gateway.HttpService
}

const (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

func NewEtcdStore(prefix string) *EtcdStore {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return &EtcdStore{}
	}
	return &EtcdStore{
		prefix:     prefix,
		HttpPrefix: fmt.Sprintf("%s/httpserver", prefix),
		rawClient:  cli,
	}
}

func (e *EtcdStore) AddHttpService(h *gateway.HttpService) error {
	content, err := proto.Marshal(h)
	if err != nil {
		logx.Info("proto.Marshal err:%s", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = e.rawClient.Put(ctx, e.getApiKey(h), string(content))
	cancel()
	if err != nil {
		logx.Info("etch put err:%s", err.Error())
		return err
	}
	return nil
}

func (e *EtcdStore) GetHttpServices() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := e.rawClient.Get(ctx, e.getHttpPrefix(), clientv3.WithPrefix())
	cancel()
	if err != nil {
		logx.Info("etch put err:%s", err.Error())
		return err
	}
	for _, ev := range resp.Kvs {
		service := &gateway.HttpService{}
		err = proto.Unmarshal(ev.Value, service)
		if err != nil {
			logx.Info("Unmarshal err:%s", err.Error())
			return err
		}
		e.HttpServices = append(e.HttpServices, service)
	}
	return nil
}

func (e *EtcdStore) getApiKey(h *gateway.HttpService) string {
	return fmt.Sprintf("%s/%d", e.HttpPrefix, h.ServerInfo.Id)
}

func (e *EtcdStore) getHttpPrefix() string {
	return e.HttpPrefix
}

