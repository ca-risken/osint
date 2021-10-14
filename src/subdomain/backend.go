package main

import (
	"github.com/gassara-kys/envconfig"
)

type backendConfig struct {
	//	Port     string `default:"19001"`
	LogLevel     string `split_words:"true" default:"debug"`
	HarvesterDir string `split_words:"true" default:"/theHarvester"`
}

func newBackendConfig() (*backendConfig, error) {
	config := &backendConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	return config, nil
}
