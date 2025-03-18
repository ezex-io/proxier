package config

import (
	"errors"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server *ServerConfig `yaml:"server"`
	Proxy  []*ProxyRule  `yaml:"proxy"`
}

type ServerConfig struct {
	Host       string `yaml:"host"`
	ListenPort string `yaml:"listen_port"`
}

type ProxyRule struct {
	Endpoint       string `yaml:"endpoint"`
	DestinationURL string `yaml:"destination_url"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if err := config.basicCheck(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) basicCheck() error {
	if c.Server == nil {
		return errors.New("server configuration is missing")
	}

	if c.Server.Host == "" {
		return errors.New("server.host cannot be empty")
	}
	if c.Server.ListenPort == "" {
		return errors.New("server.listen_port cannot be empty")
	}

	if len(c.Proxy) == 0 {
		return errors.New("at least one proxy rule must be defined")
	}

	seenEndpoints := make(map[string]bool)

	for _, rule := range c.Proxy {
		if rule.Endpoint == "" {
			return errors.New("proxy rule endpoint cannot be empty")
		}
		if rule.DestinationURL == "" {
			return errors.New("proxy rule destination_url cannot be empty")
		}

		if seenEndpoints[rule.Endpoint] {
			return errors.New("duplicate proxy endpoint: " + rule.Endpoint)
		}
		seenEndpoints[rule.Endpoint] = true

		if _, err := url.ParseRequestURI(rule.DestinationURL); err != nil {
			return errors.New("invalid URL in proxy rule: " + rule.DestinationURL)
		}
	}

	return nil
}
