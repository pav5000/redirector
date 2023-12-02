package main

import (
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Verbose   int              `yaml:"verbose"`
	Redirects []SingleRedirect `yaml:"redirects"`
}

type SingleRedirect struct {
	Src     string `yaml:"src"`
	Dst     string `yaml:"dst"`
	UnixDst string `yaml:"unix-dst"`
}

func parseConfig() (Config, error) {
	rawData, err := os.ReadFile("config.yml")
	if err != nil {
		return Config{}, errors.Wrap(err, "os.ReadFile")
	}

	var conf Config
	err = yaml.Unmarshal(rawData, &conf)
	if err != nil {
		return Config{}, errors.Wrap(err, "yaml.Unmarshal")
	}

	return conf, nil
}
