package main

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	FailFast bool    `yaml:"fail_fast,omitempty"`
	Route    []Route `yaml:"route"`
}

type Route struct {
	Name   string        `yaml:"name"`
	IP     string        `yaml:"ip"`
	Expect []RouteExpect `yaml:"expect"`
}

type RouteExpect struct {
	When RouteExpectWhen `yaml:"when"`
}

type RouteExpectWhen struct {
	Device  string `yaml:"device,omitempty"`
	Gateway string `yaml:"gateway,omitempty"`
	Source  string `yaml:"source,omitempty"`
}

func ReadConfig(file string) (Config, error) {
	var config Config

	buf, err := os.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("cannot read file")
	}

	decoder := yaml.NewDecoder(bytes.NewReader(buf))
	decoder.KnownFields(true)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf(
			"could not parse config file=\"%s\" %s", file, err.Error(),
		)
	}

	return config, nil
}
