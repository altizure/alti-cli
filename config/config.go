package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jackytck/alti-cli/errors"
	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

// Config stores the config of accessing an api server.
type Config struct {
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Token    string `yaml:"token"`
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
	dir := fmt.Sprintf("%s/.altizure", home)
	if _, err2 := os.Stat(dir); os.IsNotExist(err2) {
		os.Mkdir(dir, 0755)
	}
	err = ioutil.WriteFile(dir+"/config", data, 0644)
	return err
}

// DefaultConfig returns the default endpoint and api key.
func DefaultConfig() Config {
	return Config{
		Endpoint: DefaultEndpoint,
		Key:      DefaultAppKey,
	}
}

// Load loads config from default path.
func Load() (Config, error) {
	dc := DefaultConfig()
	home, err := homedir.Dir()
	if err != nil {
		return dc, err
	}
	config := fmt.Sprintf("%s/.altizure/config", home)
	data, err := ioutil.ReadFile(config)
	if err != nil {
		return dc, errors.ErrNoConfig
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return dc, err
	}
	return c, nil
}
