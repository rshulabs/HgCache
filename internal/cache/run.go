package cache

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rshulabs/HgCache/internal/cache/config"
	"github.com/rshulabs/HgCache/internal/cache/controller"
	"github.com/rshulabs/HgCache/internal/cache/impl"
	v2 "github.com/rshulabs/HgCache/internal/cache/v2"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 运行函数
func Run(cfg *config.Config) error {
	fmt.Println(cfg.String())
	_ = createGroup("test")
	srv := impl.NewGrpcService(cfg)
	http := controller.NewHttpService(cfg)
	if cfg.App.IsStartHttp {
		go http.Start()
	}
	go srv.Start()
	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	select {
	case <-ch:
		http.WithStop()
		if err := srv.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func createGroup(name string) *v2.Group {
	return v2.NewGroup(name, 2<<10, v2.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}
