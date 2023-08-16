package impl

import (
	"context"
	"fmt"

	"github.com/rshulabs/HgCache/internal/cache/impl/pb"
	"github.com/rshulabs/HgCache/pkg/logx"
)

var db = map[string]string{
	"yyy": "xxx",
}

func (s *Impl) Get(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	logx.Info("远程调用成功")
	key := in.GetKey()
	if v, ok := db[key]; ok {
		return &pb.Response{Value: []byte(v)}, nil
	}
	return nil, fmt.Errorf("key %v not found", key)
}
