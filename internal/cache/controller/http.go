package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rshulabs/HgCache/internal/cache/config"

	"github.com/rshulabs/HgCache/pkg/logx"
)

type HttpService struct {
	server *http.Server
	r      gin.IRouter
	addr   string
}

func NewHttpService(cfg *config.Config) *HttpService {
	r := gin.Default()
	r = installHttpRouter(r)
	srv := &http.Server{
		ReadTimeout:       60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Addr:              cfg.Http.Addr(),
		Handler:           r,
	}
	return &HttpService{
		server: srv,
		r:      r,
		addr:   cfg.Http.Addr(),
	}
}

func (s *HttpService) Start() error {
	logx.Infof("http 服务监听地址: %s", s.addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logx.Info("server closed successfully")
			return nil
		}
		return err
	}
	return nil
}

func (s *HttpService) WithStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
