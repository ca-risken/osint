package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"time"

	"github.com/gassara-kys/envconfig"
)

type websiteClient struct {
	config websiteConfig
}

type websiteConfig struct {
	WappalyzerPath string `required:"true" split_words:"true" default:"/opt/wappalyzer/src/drivers/npm/cli.js"`
}

func newWappalyzerClient() (websiteClient, error) {
	var conf websiteConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		return websiteClient{}, err
	}
	cli := websiteClient{
		config: conf,
	}
	return cli, nil
}

func (c *websiteClient) run(target string) (*wappalyzerResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "node", c.config.WappalyzerPath, target, "-r")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		appLogger.Errorf("Failed to execute wappalyzer. error: %v, stderr: %v", err, stderr.String())
		return nil, err
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
