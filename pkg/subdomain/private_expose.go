package subdomain

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/finding"
	"github.com/vikyd/zero"
)

func searchPrivateExpose(host host, detectList *[]string, logger logging.Logger) privateExpose {
	if !zero.IsZeroVal(host.IP) && !zero.IsZeroVal(host.HostName) {
		http, urlHTTP := getHTTPStatus(host.HostName, "http", logger)
		https, urlHTTPS := getHTTPStatus(host.HostName, "https", logger)
		isDetect := isDetected(host.HostName, detectList)
		if !zero.IsZeroVal(http) && !zero.IsZeroVal(https) {
			return privateExpose{HostName: host.HostName, HTTP: http, URLHTTP: urlHTTP, HTTPS: https, URLHTTPS: urlHTTPS, IsDetected: isDetect}
		}
	}
	return privateExpose{}
}

func getHTTPStatus(host, protocol string, logger logging.Logger) (int, string) {
	res := requestHTTP(host, protocol, logger)
	if res == nil {
		return 0, ""
	}
	defer res.Body.Close()
	return res.StatusCode, res.Request.URL.String()
}

func isDetected(host string, detectList *[]string) bool {
	for _, detectWord := range *detectList {
		if strings.Contains(host, detectWord) {
			return true
		}
	}
	return false
}

func (p *privateExpose) makeFinding(projectID uint32, dataSource string) (*finding.FindingForUpsert, error) {
	if zero.IsZeroVal(*p) || !p.IsDetected {
		return nil, nil
	}
	data, err := json.Marshal(map[string]privateExpose{"data": *p})
	if err != nil {
		return nil, err
	}
	finding := &finding.FindingForUpsert{
		Description:      p.getDescription(),
		DataSource:       dataSource,
		DataSourceId:     generateDataSourceID(fmt.Sprintf("private_expose_%v", p.HostName)),
		ResourceName:     p.HostName,
		ProjectId:        projectID,
		OriginalScore:    p.getScore(),
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
	if len(desc) > 200 {
		desc = desc[:196] + " ..." // cut long text
	}
	return desc
}

type privateExpose struct {
	HostName   string `json:"hostname"`
	HTTP       int    `json:"http"`
	URLHTTP    string `json:"url_http"`
	HTTPS      int    `json:"https"`
	URLHTTPS   string `json:"url_https"`
	IsDetected bool   `json:"-"`
}
