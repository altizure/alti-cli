package config

import (
	"io/ioutil"
	"path"

	"github.com/jackytck/alti-cli/errors"
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

// GetConfigDir returns the path of the config directory that
// contains the config.yaml file.
func GetConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".altizure"), nil
}

// Config represents everything in the config stored by viper.
type Config struct {
	Scopes map[string]Scope `yaml:"scopes"`
	Active string           `yaml:"active"` // active profile id
}

// GetActive returns the active endpoint and profile of current config.
func (c Config) GetActive() APoint {
	ret := APoint{
		Name:     "Default",
		Endpoint: DefaultEndpoint,
		Key:      DefaultAppKey,
	}
	for _, v := range c.Scopes {
		for _, p := range v.Profiles {
			if p.ID == c.Active {
				ret.Name = p.Name
				ret.Endpoint = v.Endpoint
				ret.Key = p.Key
				ret.Token = p.Token
				return ret
			}
		}
	}
	return ret
}

// GetProfile finds the closest profile that matches the given id.
// If no match, return ErrProfileNotFound
func (c Config) GetProfile(id string) (*Profile, error) {
	var ret Profile
	found := false
	score := 0
	for _, v := range c.Scopes {
		for _, p := range v.Profiles {
			s := compareStrs(id, p.ID)
			if s > score {
				found = true
				score = s
				ret = p
			}
		}
	}
	if !found {
		return nil, errors.ErrProfileNotFound
	}
	return &ret, nil
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

// RemoveProfile removes the closest profile that matches the given id.
// If no match, return ErrProfileNotFound
// If active profile is removed, default profile will be set as active.
// Default profile could not be removed.
// Return the removed APoint if succeed.
func (c Config) RemoveProfile(id string, save bool) (*APoint, error) {
	var ret APoint

	// a. find closest profile
	var scope string
	var profile Profile
	var pIndex int
	found := false
	score := 0
	for k, v := range c.Scopes {
		for i, p := range v.Profiles {
			s := compareStrs(id, p.ID)
			if s > score {
				found = true
				scope = k
				score = s
				profile = p
				pIndex = i
				ret.Endpoint = v.Endpoint
				ret.Key = p.Key
				ret.Name = p.Name
				ret.Email = p.Email
				ret.Token = p.Token
			}
		}
	}
	if !found {
		return nil, errors.ErrProfileNotFound
	}
	if profile.ID == DefaultProfileID {
		return nil, errors.ErrProfileNotRemovable
	}

	// b. set new scope with removed profile
	pSlice := c.Scopes[scope].Profiles
	c.Scopes[scope] = Scope{
		Endpoint: c.Scopes[scope].Endpoint,
		Profiles: append(pSlice[:pIndex], pSlice[pIndex+1:]...),
	}

	// c. set default profile as active if active profile is removed
	if profile.ID == c.Active {
		c.Active = DefaultProfileID
	}

	// d. save
	if save {
		err := c.Save()
		if err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// ClearActiveToken clears the token of active profile.
func (c *Config) ClearActiveToken(save bool) error {
	var s string
	for k, v := range c.Scopes {
		for i, p := range v.Profiles {
			if p.ID == c.Active {
				s = k
				v.Profiles[i].Token = ""
			}
		}
	}
	c.Scopes[s] = Scope{
		Endpoint: c.Scopes[s].Endpoint,
		Profiles: uniqueProfile(c.Scopes[s].Profiles),
	}
	if c.Size() == 1 {
		c.Active = "default"
	}
	if save {
		return c.Save()
	}
	return nil
}

// Size counts the number of profiles.
func (c *Config) Size() int {
	var ret int
	for _, s := range c.Scopes {
		ret += len(s.Profiles)
	}
	return ret
}

// SetActiveUserInfo sets the username of the active profile.
func (c *Config) SetActiveUserInfo(name, email string, save bool) error {
	for _, v := range c.Scopes {
		for i, p := range v.Profiles {
			if p.ID == c.Active {
				v.Profiles[i].Name = name
				v.Profiles[i].Email = email
			}
		}
	}
	if save {
		return c.Save()
	}
	return nil
}

// Save saves the config in default path: '~/.altizure/config'.
func (c Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	confDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	confPath := path.Join(confDir, "config.yaml")
	err = ioutil.WriteFile(confPath, data, 0644)
	if err != nil {
		return err
	}
	return viper.ReadInConfig()
}

func (c Config) String() string {
	data, _ := yaml.Marshal(c)
	return string(data)
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
		if v.Equal(p) {
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
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
	Token string `yaml:"token"`
}

// Equal commpares if two profiles are equal, ignoring id.
func (p Profile) Equal(o Profile) bool {
	return p.Key == o.Key && p.Token == o.Token
}

// APoint represents the active endping and profile.
type APoint struct {
	Endpoint string `yaml:"endpoint"`
	Name     string `yaml:"name"`
	Email    string `yaml:"email"`
	Key      string `yaml:"key"`
	Token    string `yaml:"token"`
}
