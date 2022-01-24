package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/proto/finding"
)

func TestGetTakeOverScore(t *testing.T) {
	cases := []struct {
		name         string
		baseTakeover takeover
		isDown       bool
		want         float32
	}{
		{
			name: "Domain doesn't match list. Server is down.",
			baseTakeover: takeover{
				Domain: "hogehogedomain.com",
				CName:  "cname.com",
			},
			isDown: true,
			want:   6.0,
		},
		{
			name: "Domain matches list. Server is down.",
			baseTakeover: takeover{
				Domain: "hogehogedomain.com",
				CName:  "cname.github.io",
			},
			isDown: true,
			want:   8.0,
		},
		{
			name: "Domain matches list. Server is up.",
			baseTakeover: takeover{
				Domain: "hogehogedomain.com",
				CName:  "cname.github.io",
			},
			isDown: false,
			want:   1.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.baseTakeover.getScore(c.isDown)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestGetTakeOverDescription(t *testing.T) {
	cases := []struct {
		name         string
		baseTakeover takeover
		isDown       bool
		want         string
	}{
		{
			name: "Server is down.",
			baseTakeover: takeover{
				Domain: "hogehogedomain.com",
				CName:  "cname.com",
			},
			isDown: true,
			want:   "hogehogedomain.com seems to be down. It has subdomain takeover risk.(CName: cname.com)",
		},
		{
			name: "Server is up.",
			baseTakeover: takeover{
				Domain: "hogehogedomain.com",
				CName:  "cname.github.io",
			},
			isDown: false,
			want:   "hogehogedomain.com has a CName record.(CName: cname.github.io)",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.baseTakeover.getDescription(c.isDown)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestMakeTakeoverFinding(t *testing.T) {
	cases := []struct {
		name         string
		baseTakeover takeover
		isDown       bool
		projectID    uint32
		dataSource   string
		want         *finding.FindingForUpsert
		wantErr      bool
	}{
		{
			name: "Success",
			baseTakeover: takeover{
				Domain: "hogehogedomain.com",
				CName:  "cname.com",
			},
			isDown:     true,
			projectID:  1,
			dataSource: "dataSource",
			want: &finding.FindingForUpsert{
				Description:      "",
				DataSource:       "dataSource",
				DataSourceId:     generateDataSourceID("hogehogedomain.com_cname.com"),
				ResourceName:     "hogehogedomain.com",
				ProjectId:        1,
				OriginalScore:    0.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			wantErr: false,
		},
		{
			name:         "Blank takeover",
			baseTakeover: takeover{},
			isDown:       true,
			projectID:    1,
			dataSource:   "dataSource",
			want:         nil,
			wantErr:      false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want != nil {
				c.want.Description = c.baseTakeover.getDescription(c.isDown)
				c.want.OriginalScore = c.baseTakeover.getScore(c.isDown)
				data, _ := json.Marshal(map[string]takeover{"data": c.baseTakeover})
				c.want.Data = string(data)
			}
			got, err := c.baseTakeover.makeFinding(c.isDown, c.projectID, c.dataSource)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: \nwant=%v, \n got=%v", c.want, got)
			}
			if c.wantErr == (err == nil) {
				t.Fatalf("Unexpected error: wantError=%v, gotError=%v", c.wantErr, err)
			}
		})
	}
}

func TestGetPrivateExposeScore(t *testing.T) {
	cases := []struct {
		name              string
		basePrivateExpose privateExpose
		want              float32
	}{
		{
			name: "Accessible",
			basePrivateExpose: privateExpose{
				HostName:   "hogehoge.com",
				HTTP:       200,
				URLHTTP:    "http://hogehoge.com",
				HTTPS:      200,
				URLHTTPS:   "https://hogehoge.com",
				IsDetected: true,
			},
			want: 6.0,
		},
		{
			name: "Not accessible",
			basePrivateExpose: privateExpose{
				HostName:   "hogehoge.com",
				HTTP:       403,
				URLHTTP:    "http://hogehoge.com",
				HTTPS:      403,
				URLHTTPS:   "https://hogehoge.com",
				IsDetected: true,
			},
			want: 3.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.basePrivateExpose.getScore()
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestGetPrivateExposeDescription(t *testing.T) {
	cases := []struct {
		name              string
		basePrivateExpose privateExpose
		want              string
	}{
		{
			name: "Accessible (http/https)",
			basePrivateExpose: privateExpose{
				HostName:   "hogehoge.com",
				HTTP:       200,
				URLHTTP:    "http://hogehoge.com",
				HTTPS:      200,
				URLHTTPS:   "https://hogehoge.com",
				IsDetected: true,
			},
			want: "hogehoge.com is accessible from public. (http/https)",
		},
		{
			name: "Accessible (https)",
			basePrivateExpose: privateExpose{
				HostName:   "hogehoge.com",
				HTTP:       0,
				URLHTTP:    "http://hogehoge.com",
				HTTPS:      200,
				URLHTTPS:   "https://hogehoge.com",
				IsDetected: true,
			},
			want: "hogehoge.com is accessible from public. (https)",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.basePrivateExpose.getDescription()
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestMakePrivateExposeFinding(t *testing.T) {
	cases := []struct {
		name              string
		basePrivateExpose privateExpose
		projectID         uint32
		dataSource        string
		want              *finding.FindingForUpsert
		wantErr           bool
	}{
		{
			name: "Success",
			basePrivateExpose: privateExpose{
				HostName:   "hogehoge.com",
				HTTP:       200,
				URLHTTP:    "http://hogehoge.com",
				HTTPS:      200,
				URLHTTPS:   "https://hogehoge.com",
				IsDetected: true,
			},
			projectID:  1,
			dataSource: "dataSource",
			want: &finding.FindingForUpsert{
				Description:      "",
				DataSource:       "dataSource",
				DataSourceId:     generateDataSourceID("private_expose_hogehoge.com"),
				ResourceName:     "hogehoge.com",
				ProjectId:        1,
				OriginalScore:    0.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			wantErr: false,
		},
		{
			name:              "Blank privateExpose",
			basePrivateExpose: privateExpose{},
			dataSource:        "dataSource",
			want:              nil,
			wantErr:           false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want != nil {
				c.want.Description = c.basePrivateExpose.getDescription()
				c.want.OriginalScore = c.basePrivateExpose.getScore()
				data, _ := json.Marshal(map[string]privateExpose{"data": c.basePrivateExpose})
				c.want.Data = string(data)
			}
			got, err := c.basePrivateExpose.makeFinding(c.projectID, c.dataSource)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: \nwant=%v, \n got=%v", c.want, got)
			}
			if c.wantErr == (err == nil) {
				t.Fatalf("Unexpected error: wantError=%v, gotError=%v", c.wantErr, err)
			}
		})
	}
}

func TestGetCertificateExpirationScore(t *testing.T) {
	cases := []struct {
		name                      string
		baseCertificateExpiration certificateExpiration
		want                      float32
	}{
		{
			name: "Score 1.0",
			baseCertificateExpiration: certificateExpiration{
				URL:        "https://hogehoge.com",
				ExpireDate: time.Now().AddDate(0, 0, 30),
			},
			want: 1.0,
		},
		{
			name: "Score 6.0",
			baseCertificateExpiration: certificateExpiration{
				URL:        "https://hogehoge.com",
				ExpireDate: time.Now().AddDate(0, 0, 14),
			},
			want: 6.0,
		},
		{
			name: "Score 8.0",
			baseCertificateExpiration: certificateExpiration{
				URL:        "https://hogehoge.com",
				ExpireDate: time.Now().AddDate(0, 0, 13),
			},
			want: 8.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.baseCertificateExpiration.getScore()
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestGetCertificateExpirationDescription(t *testing.T) {
	cases := []struct {
		name                      string
		baseCertificateExpiration certificateExpiration
		want                      string
	}{
		{
			name: "OK",
			baseCertificateExpiration: certificateExpiration{
				URL:        "https://hogehoge.com",
				ExpireDate: time.Now(),
			},
			want: fmt.Sprintf("The security certificate expires on %v, url: %v", time.Now().Format("2006-01-02"), "https://hogehoge.com"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.baseCertificateExpiration.getDescription()
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestMakeCertificateExpirationFinding(t *testing.T) {
	cases := []struct {
		name                      string
		baseCertificateExpiration certificateExpiration
		projectID                 uint32
		dataSource                string
		want                      *finding.FindingForUpsert
		wantErr                   bool
	}{
		{
			name: "Success",
			baseCertificateExpiration: certificateExpiration{
				URL:        "https://hogehoge.com",
				ExpireDate: time.Now(),
			},
			projectID:  1,
			dataSource: "dataSource",
			want: &finding.FindingForUpsert{
				Description:      "",
				DataSource:       "dataSource",
				DataSourceId:     generateDataSourceID(fmt.Sprintf("https://hogehoge.com_certificate_%v", time.Now().Format("2006-01-02"))),
				ResourceName:     "https://hogehoge.com",
				ProjectId:        1,
				OriginalScore:    0.0,
				OriginalMaxScore: 10.0,
				Data:             "",
			},
			wantErr: false,
		},
		{
			name:                      "Blank privateExpose",
			baseCertificateExpiration: certificateExpiration{},
			dataSource:                "dataSource",
			want:                      nil,
			wantErr:                   false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want != nil {
				c.want.Description = c.baseCertificateExpiration.getDescription()
				c.want.OriginalScore = c.baseCertificateExpiration.getScore()
				data, _ := json.Marshal(map[string]certificateExpiration{"data": c.baseCertificateExpiration})
				c.want.Data = string(data)
			}
			got, err := c.baseCertificateExpiration.makeFinding(c.projectID, c.dataSource)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: \nwant=%v, \n got=%v", c.want, got)
			}
			if c.wantErr == (err == nil) {
				t.Fatalf("Unexpected error: wantError=%v, gotError=%v", c.wantErr, err)
			}
		})
	}
}
