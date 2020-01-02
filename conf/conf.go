package conf

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
)

var (
	// ConfigFileName is a name of botkube configuration file
	configFile = "config.yaml"
)

type Config struct {
	Settings Settings `yaml:"settings"`
}

type Settings struct {
	AutoCreateSvc     bool           `yaml:"auto-create-svc"`
	AutocreateSvcType v1.ServiceType `yaml:"auto-create-svc-type"`
}

func New() (*Config, error) {
	c := &Config{}
	file, err := os.Open(configFile)
	defer file.Close()
	if err != nil {
		return c, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return c, err
	}

	if len(b) != 0 {
		yaml.Unmarshal(b, c)
	}
	return c, nil
}
