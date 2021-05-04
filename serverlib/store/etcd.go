package store

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"kisgateway/serverlib/logx"
	"kisgateway/serverlib/store/gateway"
	"strings"
	"time"
)

const (
	EtcdEventAdd = iota
	EtcdEventDelete
	EtcdEventUpdate

	EtcdSrcHttp = iota
)

type HandlerFunc func(watch *EtcdWatch)
type EtcdStore struct {
	prefix       string
	HttpPrefix   string
	HttpServices map[string]*gateway.ServiceInfo
	rawClient    *clientv3.Client

	eventCh chan EtcdWatch
	HandleMap map[int]HandlerFunc
}

type EtcdWatch struct {
	Type int //类型增删改查
	Src int //handle索引
	Key string
	Value *gateway.ServiceInfo
}

type Services struct {
	HttpServicesM map[string]*gateway.ServiceInfo
	HttpService   *gateway.ServiceInfo
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
	etcdStore := &EtcdStore{
		prefix:     prefix,
		HttpPrefix: fmt.Sprintf("%s/httpserver", prefix),
		rawClient:  cli,
	}

	etcdStore.init()
	return etcdStore
}

func (e *EtcdStore) AddHttpService(h *gateway.ServiceInfo) error {
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
		service := &gateway.ServiceInfo{}
		err = proto.Unmarshal(ev.Value, service)
		if err != nil {
			logx.Info("Unmarshal err:%s", err.Error())
			return err
		}
		e.HttpServices[service.ServerInfo.ServiceName] = service
	}
	return nil
}

func (e *EtcdStore) Watch() {
	watcher := clientv3.NewWatcher(e.rawClient)
	defer watcher.Close()

	ctx := e.rawClient.Ctx()
	for {
		rch := watcher.Watch(ctx, e.prefix, clientv3.WithPrefix())
		for wresp := range rch {
			if wresp.Canceled {
				return
			}
			var event EtcdWatch
			event = EtcdWatch{}
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.DELETE:
					event.Type = EtcdEventDelete
				case mvccpb.PUT:
					if ev.IsCreate() {
						event.Type = EtcdEventAdd
					} else if ev.IsModify() {
						event.Type = EtcdEventUpdate
					}
				}
				key := string(ev.Kv.Key)
				if strings.HasPrefix(key, e.HttpPrefix) {
					event.Src = EtcdSrcHttp
				}
				event.Key = key
				v := gateway.ServiceInfo{}
				err := proto.UnmarshalMerge(ev.Kv.Value, &v)
				if err != nil {
					break
				}
				event.Value = &v
				e.eventCh <- event
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (e *EtcdStore)HandleEventCh()  {
	select {
	case event := <-e.eventCh:
		e.handleWatch(event)
	default:
	}
}

func (e *EtcdStore)handleWatch(event EtcdWatch)  {
	e.HandleMap[event.Src](&event)
}

func (e *EtcdStore) getApiKey(h *gateway.ServiceInfo) string {
	return fmt.Sprintf("%s/%d", e.HttpPrefix, h.ServerInfo.Id)
}

func (e *EtcdStore) getHttpPrefix() string {
	return e.HttpPrefix
}

func (e *EtcdStore) init() {
	e.HandleMap = make(map[int]HandlerFunc)
	e.HandleMap[EtcdSrcHttp] = func(watch *EtcdWatch) {
		if watch.Type == EtcdEventDelete {
			delete(e.HttpServices, watch.Key)
		}
		if watch.Type == EtcdEventUpdate || watch.Type == EtcdEventAdd {
			e.HttpServices[watch.Key] = watch.Value
		}
	}
	go e.Watch()
	go e.HandleEventCh()
}
