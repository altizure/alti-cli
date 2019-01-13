package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

// Config stores the config of accessing an api server.
type Config struct {
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Token    string `yaml:"token"`
}

// Scope returns the same key if two configs belong to the same scope.
// A scope is defined by the scheme and host of endpoint.
func (c Config) Scope() string {
	u, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return strings.Replace(DefaultScope, ".", "*", -1)
	}
	s := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	if u.Port() != "" {
		s += ":" + u.Port()
	}
	return strings.Replace(s, ".", "*", -1)
}

// Save saves the config in default path: '~/.altizure/config'.
func (c Config) Save() error {
	viper.Set("endpoint", c.Endpoint)
	viper.Set("key", c.Key)
	viper.Set("token", c.Token)
	return viper.WriteConfig()
}

// DefaultConfig returns the default endpoint and api key.
func DefaultConfig() Config {
	return Config{
		Endpoint: viper.GetString("endpoint"),
		Key:      viper.GetString("key"),
	}
}

// Load loads config from default path.
func Load() Config {
	c := Config{
		Endpoint: viper.GetString("endpoint"),
		Key:      viper.GetString("key"),
		Token:    viper.GetString("token"),
	}
	return c
}
