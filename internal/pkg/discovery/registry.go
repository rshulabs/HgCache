package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rshulabs/HgCache/pkg/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 服务类型
type Service interface {
	Name() string
	Addr() string
}

type Registry interface {
	// 注册
	Regisrty(srv Service) error
	// 注销
	DeRegistry(srv Service) error
}

// 基于etcd服务发现中间件，实现regisrty
type EtcdRegisty struct {
	// 客户端信息
	cli *clientv3.Client
	// 租约信息
	leaseID clientv3.LeaseID
	// 租约时间
	leaseTTL int64
	// 续约响应 chan
	leaseKeepAliveRespCh <-chan *clientv3.LeaseKeepAliveResponse
}

type Option func(er *EtcdRegisty)

func WithLeaseTTL(ttl int64) Option {
	return func(er *EtcdRegisty) {
		er.leaseTTL = ttl
	}
}

func NewEtcdRegisty(endpoints []string, opts ...Option) (*EtcdRegisty, error) {
	// 校验
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("endpoints cannot be empty")
	}
	er := &EtcdRegisty{}
	for _, opt := range opts {
		opt(er)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	er.cli = cli
	return er, nil
}

func (r *EtcdRegisty) Registy(srv Service) error {
	// 申请租约
	grantRes, err := r.cli.Grant(context.Background(), r.leaseTTL)
	if err != nil {
		return err
	}
	// 得到leaseID
	r.leaseID = grantRes.ID
	// 前缀key
	key := fmt.Sprintf("%s-%s", srv.Name(), uuid.New().String())
	// 将服务put到etcd
	_, err = r.cli.Put(context.Background(), key, srv.Addr(), clientv3.WithLease(r.leaseID))
	if err != nil {
		return err
	}
	// 续约
	r.leaseKeepAliveRespCh, err = r.cli.KeepAlive(context.Background(), r.leaseID)
	if err != nil {
		return err
	}
	// 并发接收续约chan
	go HandleKeepAliveResp(r.leaseKeepAliveRespCh)
	logx.Infof("%s server is registered.", srv.Name())
	return nil
}

func (r *EtcdRegisty) DeRegisty() error {
	// revoke
	if _, err := r.cli.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	// close etcd client
	if err := r.cli.Close(); err != nil {
		return err
	}
	return nil
}

func HandleKeepAliveResp(ch <-chan *clientv3.LeaseKeepAliveResponse) {
	for resp := range ch {
		logx.Infof("service is keepalived with %x", resp.ID)
	}
}
