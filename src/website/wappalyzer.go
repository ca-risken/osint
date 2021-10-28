package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/gassara-kys/envconfig"
)

type wappalyzerClient struct {
	config wappalyzerConfig
}

type wappalyzerConfig struct {
	ResultDir      string `required:"true" split_words:"true" default:"/tmp"`
	WappalyzerPath string `required:"true" split_words:"true" default:"/opt/wappalyzer/cli.js"`
	AWSRegion      string `envconfig:"aws_region"             default:"ap-northeast-1"`
}

func newWappalyzerClient() (wappalyzerClient, error) {
	var conf wappalyzerConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		return wappalyzerClient{}, err
	}
	cli := wappalyzerClient{
		config: conf,
	}
	return cli, nil
}

func (c *wappalyzerClient) run(target string) (*wappalyzerResult, error) {
	appLogger.Infof("config: %v", c.config)
	cmd := exec.Command(c.config.WappalyzerPath, target)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		appLogger.Errorf("Failed to execute wappalyzer.error: %v, stderr: %v", err, stderr.String())
		return nil, fmt.Errorf("Failed to execute wappalyzer. error: %v", err)
	}

	var result wappalyzerResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		appLogger.Errorf("Failed to parse scan result. error: %v", err)
		return nil, err
	}

	return &result, nil
}

type wappalyzerResult struct {
	URLs         interface{}
	Technologies []wappalyzerTechnology
}

type wappalyzerTechnology struct {
	Slug       string
	Name       string
	Confidence int
	Version    string
	Icon       string
	Website    string
	CPE        string
	Categories []wappalyzerCategory
}

type wappalyzerCategory struct {
	ID   int
	Slug string
	Name string
}
