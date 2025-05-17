package config

import (
	"github.com/ezex-io/gopkg/env"
)

type Config struct {
	Address        string
	EnableFastHTTP bool
	RawRules       string
	ParsedRules    []*Rules
}

type Rules struct {
	Endpoint    string `json:"endpoint"`
	Destination string `json:"destination"`
}

func LoadFromEnv() *Config {
	return &Config{
		Address:        env.GetEnv[string]("PROXIER_ADDRESS", env.WithDefault("0.0.0.0:8080")),
		EnableFastHTTP: env.GetEnv[bool]("PROXIER_ENABLE_FASTHTTP", env.WithDefault("false")),
		RawRules:       env.GetEnv[string]("PROXIER_RULES"),
		ParsedRules:    make([]*Rules, 0),
	}
}

func (*Config) BasicCheck() error {
	return nil
}
