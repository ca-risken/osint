package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"syscall"
	"time"

	"github.com/Songmu/timeout"
)

type websiteClient struct {
	WappalyzerPath string
}

func newWappalyzerClient(wappalyzerPath string) (websiteClient, error) {
	cli := websiteClient{
		WappalyzerPath: wappalyzerPath,
	}
	return cli, nil
}

func (c *websiteClient) run(target string) (*wappalyzerResult, error) {
	cmd := exec.Command("node", c.WappalyzerPath, target, "-r")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  20 * time.Minute,
		KillAfter: 5 * time.Second,
		Signal:    syscall.SIGTERM,
	}
	exitStatus, err := tio.RunContext(context.Background())
	if exitStatus.IsTimedOut() {
		appLogger.Errorf("Timeout occured when executing wappalyzer.")
		return nil, errors.New("Timeout occured when executing wappalyzer")
	}
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
