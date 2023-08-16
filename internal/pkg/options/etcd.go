package options

import "github.com/spf13/pflag"

type EtcdOption struct {
	Endpoints []string `json:"endpoints" mapstructure:"endpoints" yaml:"endpoints"`
	LeaseTTL  int64    `json:"leaseTTL" mapstructure:"leaseTTL" yaml:"leaseTTL"`
}

func NewEtcdOption() *EtcdOption {
	return &EtcdOption{
		Endpoints: []string{"127.0.0.1:2379"},
		LeaseTTL:  5,
	}
}

func (o *EtcdOption) Validate() []error {
	return nil
}

func (o *EtcdOption) AddFlags(fs *pflag.FlagSet) {
	fs.Int64Var(&o.LeaseTTL, "lease.ttl", 5, "The leaseTTL of etcd.")
}
