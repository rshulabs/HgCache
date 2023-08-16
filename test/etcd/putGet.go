package etcd

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func PutGet() {
	cli, err := clientv3.New(clientv3.Config{
		// etcd serevr host
		Endpoints:   []string{"192.168.60.34:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	// put
	pRes, err := cli.Put(context.Background(), "a1", "123456")
	if err != nil {
		panic(err)
	}
	fmt.Println(pRes)
	// get
	gRes, err := cli.Get(context.Background(), "a1")
	if err != nil {
		panic(err)
	}
	fmt.Println(gRes.Kvs[0].Key, gRes.Kvs[0].Value)
}
