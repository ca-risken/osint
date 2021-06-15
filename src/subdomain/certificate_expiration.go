package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/vikyd/zero"
)

func checkCertificateExpirationFromDomain(domain string) certificateExpiration {
	url := fmt.Sprintf("https://%v/", domain)
	certificateExpired := checkCertificateExpiration(url)
	if (certificateExpired == time.Time{}) {
		return certificateExpiration{}
	}
	return certificateExpiration{
		URL:        url,
		ExpireDate: certificateExpired,
	}
}

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
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}
	}
	expireUTCTime := resp.TLS.PeerCertificates[0].NotAfter
	return expireUTCTime
}

func (c *certificateExpiration) makeFinding(projectID uint32, dataSource, resourceName string) (*finding.FindingForUpsert, error) {
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
		ResourceName:     resourceName,
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
