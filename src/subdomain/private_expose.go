package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/vikyd/zero"
)

func searchPrivateExpose(host host, detectList *[]string) privateExpose {
	if !zero.IsZeroVal(host.IP) && !zero.IsZeroVal(host.HostName) && isDetected(host.HostName, detectList) {
		http, urlHTTP, _ := getHTTPStatus(host.HostName, "http")
		https, urlHTTPS, _ := getHTTPStatus(host.HostName, "https")
		if !zero.IsZeroVal(http) && !zero.IsZeroVal(https) {
			return privateExpose{HostName: host.HostName, HTTP: http, URLHTTP: urlHTTP, HTTPS: https, URLHTTPS: urlHTTPS}
		}
	}
	return privateExpose{}
}

func getHTTPStatus(host, protocol string) (int, string, error) {
	url := fmt.Sprintf("%s://%s", protocol, host)
	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)

	// Timeoutもエラーに入るので、特にログも出さないでスルー(ドメインを見つけてもHTTPで使われているとは限らないため)
	if err != nil {
		return 0, "", err
	}

	defer res.Body.Close()
	return res.StatusCode, res.Request.URL.String(), nil
}

func isDetected(host string, detectList *[]string) bool {
	for _, detectWord := range *detectList {
		if strings.Index(host, detectWord) > -1 {
			return true
		}
	}
	return false
}

func (p *privateExpose) makeFinding(domain string, projectID uint32, dataSource, resourceName string) (*finding.FindingForUpsert, error) {
	if zero.IsZeroVal(*p) {
		return nil, nil
	}
	score := p.getScore()
	description := p.getDescription()
	data, err := json.Marshal(map[string]privateExpose{"data": *p})
	if err != nil {
		return nil, err
	}
	finding := &finding.FindingForUpsert{
		Description:      description,
		DataSource:       dataSource,
		DataSourceId:     domain,
		ResourceName:     resourceName,
		ProjectId:        projectID,
		OriginalScore:    score,
		OriginalMaxScore: 10.0,
		Data:             string(data),
	}
	return finding, nil
}

func (p *privateExpose) getScore() float32 {
	var score float32 = 3.0
	if (p.HTTP != 401 && p.HTTP != 403) || (p.HTTPS != 401 && p.HTTPS != 403) {
		score = score + 3.0
	}
	return score
}

func (p *privateExpose) getDescription() string {
	desc := fmt.Sprintf("%s is accessible from public.", p.HostName)
	if !zero.IsZeroVal(p.HTTP) && !zero.IsZeroVal(p.HTTPS) {
		desc = desc + " (http/https)"
	} else if !zero.IsZeroVal(p.HTTP) {
		desc = desc + " (http)"
	} else {
		desc = desc + " (https)"
	}
	return desc
}

type privateExpose struct {
	HostName string `json:"hostname"`
	HTTP     int    `json:"http"`
	URLHTTP  string `json:"url_http"`
	HTTPS    int    `json:"https"`
	URLHTTPS string `json:"url_https"`
}
