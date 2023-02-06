package app

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/weaveworks/common/server"
	yaml "gopkg.in/yaml.v2"

	"github.com/zachfi/weigh/modules/exporter"
	"github.com/zachfi/zkit/pkg/tracing"
)

type Config struct {
	Target  string         `yaml:"target"`
	Tracing tracing.Config `yaml:"tracing"`

	// modules
	Server server.Config `yaml:"server,omitempty"`

	Exporter exporter.Config `yaml:"exporter"`
}

// LoadConfig receives a file path for a configuration to load.
func LoadConfig(file string) (Config, error) {
	filename, _ := filepath.Abs(file)

	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return config, errors.Wrap(err, "failed to load yaml file")
	}

	return config, nil
}

// loadYamlFile unmarshals a YAML file into the received interface{} or returns an error.
func loadYamlFile(filename string, d interface{}) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, d)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) RegisterFlagsAndApplyDefaults(prefix string, f *flag.FlagSet) {
	c.Target = Once
	f.StringVar(&c.Target, "target", Once, "target module")
}
