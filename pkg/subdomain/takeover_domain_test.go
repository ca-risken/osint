package subdomain

import (
	"net/http"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/jarcoal/httpmock"
)

const (
	ZETTAI_SONZAI_SHINAI = "zettai-sonzai-shinai"
)

func TestIsDownVHODomain(t *testing.T) {
	logger := logging.NewLogger()
	cases := []struct {
		name     string
		input    string
		want     bool
		mockResp string
	}{
		{
			name:     "AWS/Elastic Beanstalk",
			input:    ZETTAI_SONZAI_SHINAI + ".us-east-1.elasticbeanstalk.com",
			want:     true,
			mockResp: "",
		},
		{
			name:     "AWS/S3",
			input:    ZETTAI_SONZAI_SHINAI + ".s3-website-ap-northeast-1.amazonaws.com",
			want:     true,
			mockResp: "NoSuchBucket",
		},
		{
			name:     "Anima",
			input:    ZETTAI_SONZAI_SHINAI + ".animaapp.io",
			want:     true,
			mockResp: "Anima - Page Not Found",
		},
		{
			name:     "Bitbucket",
			input:    ZETTAI_SONZAI_SHINAI + ".bitbucket.io",
			want:     true,
			mockResp: "Repository not found",
		},
		{
			name:     "Gemfury",
			input:    ZETTAI_SONZAI_SHINAI + ".furyns.com",
			want:     true,
			mockResp: "404: This page could not be found.",
		},
		{
			name:     "Ghost",
			input:    ZETTAI_SONZAI_SHINAI + ".ghost.io",
			want:     true,
			mockResp: "Domain error",
		},
		{
			name:     "HatenaBlog",
			input:    ZETTAI_SONZAI_SHINAI + ".hatenablog.com",
			want:     true,
			mockResp: "404 Blog is not found",
		},
		{
			name:     "Help Juice",
			input:    ZETTAI_SONZAI_SHINAI + ".helpjuice.com",
			want:     true,
			mockResp: "We could not find what you're looking for.",
		},
		{
			name:     "Help Scout",
			input:    ZETTAI_SONZAI_SHINAI + ".helpscoutdocs.com",
			want:     true,
			mockResp: "No settings were found for this company:",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", c.input,
				func(req *http.Request) (*http.Response, error) {
					res := httpmock.NewStringResponse(200, c.mockResp)
					return res, nil
				},
			)
			td := getTakeoverDomain(c.input)
			if td == nil {
				t.Fatalf("Could not get takeover domain, input=%s", c.input)
			}
			got := isDownVHODomain(c.input, td.Fingerprint, logger)
			if got != c.want {
				t.Fatalf("Unexpected return: got=%t, want=%t", got, c.want)
			}
		})
	}
}
