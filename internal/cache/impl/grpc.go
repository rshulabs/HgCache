package impl

import (
	"fmt"
	"net"

	"github.com/rshulabs/HgCache/internal/cache/config"
	"github.com/rshulabs/HgCache/internal/cache/impl/pb"
	"github.com/rshulabs/HgCache/internal/pkg/discovery"
	"github.com/rshulabs/HgCache/pkg/logx"
	"google.golang.org/grpc"
)

var (
	srv = &Impl{}
)

type GrpcService struct {
	srv *grpc.Server
	cfg *config.Config
	s   *Service
	reg *discovery.EtcdRegisty
}

func NewGrpcService(cfg *config.Config) *GrpcService {
	gs := grpc.NewServer()
	pb.RegisterGroupCacheServer(gs, srv)
	s := &Service{
		SName: "cache",
		SAddr: fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port),
	}
	return &GrpcService{
		srv: gs,
		cfg: cfg,
		s:   s,
	}
}

type Service struct {
	SName string
	SAddr string
}

func (s *Service) Name() string { return s.SName }

func (s *Service) Addr() string { return s.SAddr }

func (s *GrpcService) Start() {
	// 服务注册
	reg, err := discovery.NewEtcdRegisty(s.cfg.Etcd.Endpoints, discovery.WithLeaseTTL(s.cfg.Etcd.LeaseTTL))
	if err != nil {
		panic(err)
	}
	s.reg = reg
	if err := reg.Registy(s.s); err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", s.s.SAddr)
	if err != nil {
		panic(err)
	}
	logx.Infof("GRPC 服务监听地址: %s", s.s.SAddr)
	if err := s.srv.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			logx.Info("service stopped.")
		}
		panic(err)
	}
}

func (s *GrpcService) Stop() error {
	s.reg.DeRegisty()
	s.srv.GracefulStop()
	return nil
}
