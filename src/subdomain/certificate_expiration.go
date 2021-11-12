package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ca-risken/core/proto/finding"
	"github.com/vikyd/zero"
)

func (p *privateExpose) checkCertificateExpiration() certificateExpiration {
	if p.HTTPS == 0 || p.URLHTTPS == "" {
		return certificateExpiration{}
	}
	certificateExpired := checkCertificateExpiration(p.URLHTTPS)
	if (certificateExpired == time.Time{}) {
		return certificateExpiration{}
	}
	return certificateExpiration{
		URL:        p.URLHTTPS,
		ExpireDate: certificateExpired,
	}
}

func checkCertificateExpiration(url string) time.Time {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}
	}
	if resp != nil && resp.TLS != nil && resp.TLS.PeerCertificates[0] != nil {
		return resp.TLS.PeerCertificates[0].NotAfter
	}
	return time.Time{}
}

func (c *certificateExpiration) makeFinding(projectID uint32, dataSource string) (*finding.FindingForUpsert, error) {
	if zero.IsZeroVal(*c) {
		return nil, nil
	}
	score := c.getScore()
	description := c.getDescription()
	data, err := json.Marshal(map[string]certificateExpiration{"data": *c})
	if err != nil {
		return nil, err
	}
	finding := &finding.FindingForUpsert{
		Description:      description,
		DataSource:       dataSource,
		DataSourceId:     generateDataSourceID(fmt.Sprintf("%v_%v", c.URL, "certificate")),
		ResourceName:     c.URL,
		ProjectId:        projectID,
		OriginalScore:    score,
		OriginalMaxScore: 10.0,
		Data:             string(data),
	}
	return finding, nil
}

func (c *certificateExpiration) getScore() float32 {
	now := time.Now()
	dateHighScore := now.AddDate(0, 0, 14)
	dateMiddleScore := now.AddDate(0, 0, 30)
	if c.ExpireDate.Unix() < dateHighScore.Unix() {
		return 8.0
	}
	if c.ExpireDate.Unix() < dateMiddleScore.Unix() {
		return 6.0
	}
	return 1.0
}

func (c *certificateExpiration) getDescription() string {
	expireDate := c.ExpireDate.Format("2006-01-02")
	description := fmt.Sprintf("The security certificate expires on %v, url: %v", expireDate, c.URL)
	if len(description) > 200 {
		description = description[:196] + " ..." // cut long text
	}
	return description
}

type certificateExpiration struct {
	URL        string    `json:"url"`
	ExpireDate time.Time `json:"expire_date"`
}
