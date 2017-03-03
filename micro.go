package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DataBaseConfig struct {
	Url, Name, Collection string
}

type Config struct {
	DB   DataBaseConfig `yaml:"database"`
	Port string         `yml:"port"`
}

// TODO: consider using toml files
// see https://goo.gl/dC9YVl & https://goo.gl/nNIieS
const (
	MicroFile = "micro.yml"
)

// ReadConfig reads the config file
func ReadConfig() (Config, error) {
	dat, err := ioutil.ReadFile(MicroFile)

	if err != nil {
		return Config{}, err
	}

	conf := Config{}
	err = yaml.Unmarshal(dat, &conf)
	if err != nil {
		return Config{}, err
	}

	if checkMissingValues(conf) {
		return conf, fmt.Errorf("missing values on config file %v", MicroFile)
	}

	addDefaultValues(&conf)

	return conf, nil
}

// checkMissingValues checks if there some required
// values where not provided to micro.yml
func checkMissingValues(conf Config) bool {
	return conf.DB.Url == "" || conf.DB.Name == "" || conf.DB.Collection == ""
}

// addDefaultValues adds default values if some
// defaultable fields are missing
func addDefaultValues(conf *Config) {
	if conf.Port == "" {
		conf.Port = "8080"
	}
}
