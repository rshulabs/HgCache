package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/rshulabs/HgCache/internal/pkg/discovery/lb"
	"github.com/rshulabs/HgCache/pkg/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Discovery interface {
	// 获取服务器的一个地址
	GetServiceAddr(serviceName string) (string, error)
	// 监控服务地址变化
	WatchService(serviceName string) error
}

type EtcdDiscovery struct {
	cli *clientv3.Client
}

func NewEtcdDiscovery(endpoints []string) (*EtcdDiscovery, error) {
	// 校验
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("endpoints cannot be empty")
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdDiscovery{
		cli: cli,
	}, nil
}

func (d *EtcdDiscovery) GetServiceAddr(serviceName, key string) (string, error) {
	// get --prefix
	gResp, err := d.cli.Get(context.Background(), serviceName, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(gResp.Kvs) == 0 {
		return "", fmt.Errorf("%s service is not found", serviceName)
	}
	var nodes []string
	for _, v := range gResp.Kvs {
		nodes = append(nodes, string(v.Value))
	}
	cHash := lb.NewMap(100, nil)
	cHash.Add(nodes...)
	addr := cHash.Get(key)
	return addr, nil
}

func (d *EtcdDiscovery) WatchService(serviceName string) error {
	ch := d.cli.Watch(context.Background(), serviceName, clientv3.WithPrefix())
	select {
	case <-ch:
		logx.Info("service changed")
	}
	return nil
}
