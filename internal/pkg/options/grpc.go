package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type GrpcOption struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

func NewGrpcOption() *GrpcOption {
	return &GrpcOption{
		Host: "0.0.0.0",
		Port: 9091,
	}
}

func (g *GrpcOption) Validate() []error {
	var errors []error
	if g.Port < 0 || g.Port > 65535 {
		errors = append(errors, fmt.Errorf("The grpc %v must be between 0 and 65535.", g.Port))
	}
	return errors
}

func (g *GrpcOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&g.Host, "grpc.host", g.Host, "The grpc host.")
	fs.IntVar(&g.Port, "grpc.port", g.Port, "The grpc port.")
}
