package subdomain

import (
	"testing"
)

func TestIsNXDomain(t *testing.T) {
	cases := []struct {
		domain string
		want   bool
	}{
		{
			domain: "www.example.com",
			want:   false,
		},
		{
			domain: "nonexistent.example.com",
			want:   true,
		},
		{
			domain: "security-hub.jp",
			want:   false,
		},
		{
			domain: "nonexistent.security-hub.jp",
			want:   true,
		},
	}

	for _, c := range cases {
		got := isNXDomain(c.domain)
		if got != c.want {
			t.Fatalf("Unexpected return: %v, want: %v", got, c.want)
		}
	}
}
