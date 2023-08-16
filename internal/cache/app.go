package cache

import (
	"github.com/rshulabs/HgCache/internal/cache/config"
	"github.com/rshulabs/HgCache/internal/cache/options"
	"github.com/rshulabs/micro-frame/pkg/app"
)

const commandDesc = `demo program`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	app := app.NewApp("demo api server", basename, app.WithDefaultValidArgs(), app.WithOptions(opts), app.WithDescription(commandDesc), app.WithRunFunc(run(opts)))
	return app
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		// TODO 日志初始化
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		return Run(cfg)
	}
}
