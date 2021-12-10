package main

import (
	"reflect"
	"testing"
)

func TestGetRecommend(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  recommend
	}{
		{
			name:  "Exists plugin",
			input: "Takeover",
			want: recommend{
				Type: "Domain/TakeOver",
				Risk: `Domain Takeover Lisk
		- The target domain has a CNAME.
		- If it is down or not under your control, it may be hijacked.`,
				Recommendation: `Delete unused DNS records.
		- https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/02-Configuration_and_Deployment_Management_Testing/10-Test_for_Subdomain_Takeover`,
			},
		},
		{
			name:  "Unknown category",
			input: "unknown",
			want: recommend{
				Risk:           "",
				Recommendation: "",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getRecommend(c.input)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}
