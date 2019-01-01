package config

import (
	"fmt"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

// Config stores the config of accessing an api server.
type Config struct {
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Token    string `yaml:"token"`
}

// Save saves the config in default paht.
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
	err = ioutil.WriteFile(dir+"/credentials", data, 0644)
	return err
}
