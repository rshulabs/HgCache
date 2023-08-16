package config

import "github.com/rshulabs/HgCache/internal/cache/options"

// New - viper - parser
type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
