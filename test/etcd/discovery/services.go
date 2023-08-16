package discovery

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rshulabs/HgCache/pkg/logx"
)

type UserService struct {
	UserName string
	UserAddr string
}

func NewUserService(addr string) *UserService {
	return &UserService{
		UserName: "user",
		UserAddr: addr,
	}
}

func (s UserService) Name() string {
	return s.UserName
}

func (s UserService) Addr() string {
	return s.UserAddr
}

func UserRegisty(srv *UserService) {
	ec, err := NewEtcdRegisty([]string{"192.168.60.34:2379"}, WithLeaseTTL(10))
	if err != nil {
		panic(err)
	}
	err = ec.Registy(srv)
	if err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	select {
	case <-ch:
		if err := ec.DeRegisty(); err != nil {
			logx.Fatalf(err.Error())
		}
		logx.Infof("%s is deregistryed.", srv.Name())
	}
}

func UserDiscovery(name string) (err error) {
	ed, err := NewEtcdDiscovery([]string{"192.168.60.34:2379"})
	ip, err := ed.GetServiceAddr(name)
	fmt.Println(ip)
	return
}
