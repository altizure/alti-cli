package config

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"github.com/jackytck/alti-cli/rand"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// DefaultConfig returns the default endpoint and api key.
func DefaultConfig() Config {
	m := map[string]Scope{}
	m[DefaultScope] = Scope{
		Endpoint: DefaultEndpoint,
		Profiles: []Profile{
			{
				ID:  "default",
				Key: DefaultAppKey,
			},
		},
	}
	c := Config{
		Scopes: m,
		Active: "default",
	}
	return c
}

// Load loads config from default path.
func Load() Config {
	var c Config
	err := viper.Unmarshal(&c)
	if err != nil || c.Scopes == nil {
		return DefaultConfig()
	}
	return c
}

// Config represents everything in the config stored by viper.
type Config struct {
	Scopes map[string]Scope `yaml:"scopes"`
	Active string           `yaml:"active"` // active profile id
}

// GetActive returns the active endpoint and profile of current config.
func (c Config) GetActive() APoint {
	ret := APoint{
		Endpoint: DefaultEndpoint,
		Key:      DefaultAppKey,
	}
	for _, v := range c.Scopes {
		for _, p := range v.Profiles {
			if p.ID == c.Active {
				ret.Endpoint = v.Endpoint
				ret.Key = p.Key
				ret.Token = p.Token
				return ret
			}
		}
	}
	return ret
}

// AddProfile adds a profile under its endpoint and set it as active.
// Existing values would be replaced.
func (c *Config) AddProfile(ap APoint) error {
	k := endpointToKey(ap.Endpoint)
	s, ok := c.Scopes[k]
	uid, err := rand.RememberToken()
	if err != nil {
		return err
	}
	p := Profile{
		ID:    uid,
		Key:   ap.Key,
		Token: ap.Token,
	}
	if ok {
		// scope already exists
		p = s.Add(p)
		c.Scopes[k] = s
	} else {
		// new scope
		c.Scopes[k] = Scope{
			Endpoint: ap.Endpoint,
			Profiles: []Profile{p},
		}
	}
	c.Active = p.ID
	return nil
}

// ClearActiveToken clears the token of active profile.
func (c *Config) ClearActiveToken() {
	for _, v := range c.Scopes {
		for i, p := range v.Profiles {
			if p.ID == c.Active {
				v.Profiles[i].Token = ""
			}
		}
	}
}

// Save saves the config in default path: '~/.altizure/config'.
func (c Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	confPath := path.Join(home, ".altizure", "config.yaml")
	err = ioutil.WriteFile(confPath, data, 0644)
	if err != nil {
		return err
	}
	return viper.ReadInConfig()
}

// Scope represents a certain endpint and a list of profiles.
// An endpoint is the main domain without sub-path, e.g. api.altizure.com or 127.0.0.1:8082
type Scope struct {
	Endpoint string    `yaml:"endpoint"`
	Profiles []Profile `yaml:"profiles"`
}

// Add adds a profile in this scope.
// Return the newly added or existing profile.
func (s *Scope) Add(p Profile) Profile {
	for _, v := range s.Profiles {
		if v.Key == p.Key && v.Token == p.Token {
			// already exists
			return v
		}
	}
	s.Profiles = append(s.Profiles, p)
	return p
}

// Profile represents the login profile of a user of a certain endpoint.
type Profile struct {
	ID    string `yaml:"id"`
	Key   string `yaml:"key"`
	Token string `yaml:"token"`
}

// APoint represents the active endping and profile.
type APoint struct {
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Token    string `yaml:"token"`
}

// endpointToKey returns the unique key of this scope.
// The scope is defined by the scheme, host and port of endpoint.
func endpointToKey(ep string) string {
	u, err := url.ParseRequestURI(ep)
	if err != nil {
		return strings.Replace(DefaultScope, ".", "*", -1)
	}
	str := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	if u.Port() != "" {
		str += ":" + u.Port()
	}
	return strings.Replace(str, ".", "*", -1)
}
