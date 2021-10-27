package main

import (
	"testing"
)

func TestGetScore(t *testing.T) {
	cases := []struct {
		name     string
		status   string
		category string
		plugin   string
		want     float32
	}{
		{
			name:     "OK",
			status:   "OK",
			category: "ACM",
			plugin:   "acmCertificateExpiry",
			want:     0.0,
		}, {
			name:     "WARN",
			status:   "WARN",
			category: "ACM",
			plugin:   "acmCertificateExpiry",
			want:     3.0,
		},
		{
			name:     "UNKNOWN",
			status:   "UNKNOWN",
			category: "ACM",
			plugin:   "acmCertificateExpiry",
			want:     1.0,
		},
		{
			name:     "Fail match Map",
			status:   "FAIL",
			category: "ACM",
			plugin:   "acmCertificateExpiry",
			want:     6.0,
		},
		{
			name:     "Fail not match Map",
			status:   "FAIL",
			category: "ACM",
			plugin:   "hogehogehogehoge",
			want:     3.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getScore()
			if c.want != got {
				t.Fatalf("Unexpected category name: want=%v, got=%v", c.want, got)
			}
		})
	}
}
