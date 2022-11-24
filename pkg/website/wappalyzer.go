package website

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

type WebsiteClient struct {
	WappalyzerPath string
}

func NewWappalyzerClient(wappalyzerPath string) *WebsiteClient {
	return &WebsiteClient{
		WappalyzerPath: wappalyzerPath,
	}
}

func (c *WebsiteClient) run(target string) (*wappalyzerResult, error) {
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
		return nil, errors.New("Timeout occured when executing wappalyzer")
	}
	if err != nil {
		return nil, err
	}

	var result wappalyzerResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
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
