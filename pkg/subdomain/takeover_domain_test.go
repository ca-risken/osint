package subdomain

import (
	"testing"
)

const (
	ZETTAI_SONZAI_SHINAI = "zettai-sonzai-shinai"
)

func TestIsDownVHODomain(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "AWS/Elastic Beanstalk",
			input: ZETTAI_SONZAI_SHINAI + ".us-east-1.elasticbeanstalk.com",
			want:  true,
		},
		{
			name:  "AWS/S3",
			input: ZETTAI_SONZAI_SHINAI + ".s3-website-ap-northeast-1.amazonaws.com",
			want:  true,
		},
		{
			name:  "Anima",
			input: ZETTAI_SONZAI_SHINAI + ".animaapp.io",
			want:  true,
		},
		{
			name:  "Bitbucket",
			input: ZETTAI_SONZAI_SHINAI + ".bitbucket.io",
			want:  true,
		},
		{
			name:  "Gemfury",
			input: ZETTAI_SONZAI_SHINAI + ".furyns.com",
			want:  true,
		},
		{
			name:  "Ghost",
			input: ZETTAI_SONZAI_SHINAI + ".ghost.io",
			want:  true,
		},
		{
			name:  "HatenaBlog",
			input: ZETTAI_SONZAI_SHINAI + ".hatenablog.com",
			want:  true,
		},
		{
			name:  "Help Juice",
			input: ZETTAI_SONZAI_SHINAI + ".helpjuice.com",
			want:  true,
		},
		{
			name:  "Help Scout",
			input: ZETTAI_SONZAI_SHINAI + ".helpscoutdocs.com",
			want:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			td := getTakeoverDomain(c.input)
			if td == nil {
				t.Fatalf("Could not get takeover domain, input=%s", c.input)
			}
			got := isDownVHODomain(c.input, td.Fingerprint)
			if got != c.want {
				t.Fatalf("Unexpected return: got=%t, want=%t", got, c.want)
			}
		})
	}
}
