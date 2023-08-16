package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type AppOption struct {
	Name        string `json:"name" mapstructure:"name" yaml:"name"`
	IsStartHttp bool   `json:"is_start_http" mapstructure:"is_start_http" yaml:"is_start_http"`
}

func NewAppOption() *AppOption {
	return &AppOption{
		Name:        "",
		IsStartHttp: false,
	}
}

func (a *AppOption) Validate() []error {
	var errors []error
	if a.Name == "" {
		errors = append(errors, fmt.Errorf("app name is required"))
	}
	return errors
}

func (a *AppOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&a.Name, "app.name", a.Name, "The binary name of app")
	fs.BoolVar(&a.IsStartHttp, "app.is_start_http", a.IsStartHttp, "Control the http server.")
}
