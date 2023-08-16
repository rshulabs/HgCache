package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type HttpOption struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

func NewHttpOption() *HttpOption {
	return &HttpOption{
		Host: "0.0.0.0",
		Port: 8091,
	}
}

func (h *HttpOption) Validate() []error {
	var errors []error
	if h.Port < 0 || h.Port > 65535 {
		errors = append(errors, fmt.Errorf("The http %v must be between 0 and 65535.", h.Port))
	}
	return errors
}

func (h *HttpOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&h.Host, "http.host", h.Host, "The http host.")
	fs.IntVar(&h.Port, "http.port", h.Port, "The http port.")
}

func (h *HttpOption) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}
