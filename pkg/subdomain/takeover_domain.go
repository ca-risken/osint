package subdomain

import (
	"io"
	"regexp"
	"strings"
)

type TakeoverType int

const (
	VHO TakeoverType = iota
	NVHO
)

type TakeoverDomain struct {
	ServiceName string
	Domain      *regexp.Regexp
	Type        TakeoverType
	Fingerprint string // how to check whether the subdomain has already been takeovered.
}

func isDownVHODomain(cname string, fingerprint string) bool {
	if existsVirtualHost(cname, "http", fingerprint) {
		return false
	}
	if existsVirtualHost(cname, "https", fingerprint) {
		return false
	}
	return true
}

func existsVirtualHost(cname, protocol, fingerprint string) bool {
	resp := requestHTTP(cname, protocol)
	if resp == nil || resp.Body == nil {
		return false
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(buf), fingerprint) {
		return false
	}
	for _, headers := range resp.Header {
		// Iterate all headers with one name (e.g. Content-Type)
		for _, header := range headers {
			if strings.Contains(header, fingerprint) {
				return false
			}
		}
	}
	return true
}

func getTakeoverDomain(subdomain string) *TakeoverDomain {
	for _, td := range TakeoverDomains {
		if td.Domain.MatchString(subdomain) {
			return &td
		}
	}
	return nil
}

// Domains that have takeover vulnerable status,
// source: https://github.com/EdOverflow/can-i-take-over-xyz/blob/44e2da47ecb95fc38a0976812fc173e553996189/fingerprints.json
// Cannot test domains below:
//   - agilecrm.com
//   - airee.ru
//   - youtrack.cloud
var TakeoverDomains = []TakeoverDomain{
	{
		ServiceName: "AWS/Elastic Beanstalk",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("us-east-1.elasticbeanstalk.com")),
		Type:        NVHO,
	},
	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteEndpoints.html
	{
		ServiceName: "AWS/S3",
		Domain:      regexp.MustCompile(`s3.*\.amazonaws\.com`), // fix from source
		Type:        VHO,
		Fingerprint: "NoSuchBucket", // fix from source
	},
	{
		ServiceName: "Anima",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("animaapp.io")),
		Type:        VHO,
		Fingerprint: "Anima - Page Not Found", // fix from source
	},
	{
		ServiceName: "Bitbucket",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("bitbucket.io")),
		Type:        VHO,
		Fingerprint: "Repository not found",
	},
	{
		ServiceName: "Gemfury",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("furyns.com")),
		Type:        VHO,
		Fingerprint: "404: This page could not be found.",
	},
	{
		ServiceName: "Ghost",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("ghost.io")),
		Type:        VHO,
		Fingerprint: "Domain error", // fix from source
	},
	{
		ServiceName: "HatenaBlog",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("hatenablog.com")),
		Type:        VHO,
		Fingerprint: "404 Blog is not found",
	},
	{
		ServiceName: "Help Juice",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("helpjuice.com")),
		Type:        VHO,
		Fingerprint: "We could not find what you're looking for.",
	},
	{
		ServiceName: "Help Scout",
		Domain:      regexp.MustCompile(regexp.QuoteMeta("helpscoutdocs.com")),
		Type:        VHO,
		Fingerprint: "No settings were found for this company:",
	},
}
