package main

import (
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Redirects []SingleRedirect `yaml:"redirects"`
}

type SingleRedirect struct {
	Src string `yaml:"src"`
	Dst string `yaml:"dst"`
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
