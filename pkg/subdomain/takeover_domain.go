package subdomain

import (
	"io"
	"log"
	"strings"
)

type TakeoverType int

const (
	VHO TakeoverType = iota
	NVHO
)

type TakeoverDomain struct {
	ServiceName string
	Domain      string
	Type        TakeoverType
	Fingerprint string // how to check whether the subdomain has already been takeovered.
}

func isDownVHODomain(cname string, fingerprint string) bool {
	if existsVertialHost(cname, "http", fingerprint) {
		return false
	}
	if existsVertialHost(cname, "https", fingerprint) {
		return false
	}
	return true
}

func existsVertialHost(cname, protocol, fingerprint string) bool {
	resp := requestHTTP(cname, protocol)
	if resp == nil || resp.Body == nil {
		return false
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(buf), fingerprint) {
		return false
	}
	for name, headers := range resp.Header {
		// Iterate all headers with one name (e.g. Content-Type)
		for _, header := range headers {
			log.Printf("test: cname=%s, protocol=%s, header=%s:%s", cname, protocol, name, header)
			if strings.Contains(header, fingerprint) {
				return false
			}
		}
	}
	log.Printf("test: cname=%s, protocol=%s, body=%s", cname, protocol, string(buf))
	return true
}

func getTakeoverDomain(subdomain string) *TakeoverDomain {
	for _, td := range TakeoverDomains {
		if strings.Contains(subdomain, td.Domain) {
			return &td
		}
	}
	return nil
}

// Domains that have takeover vulnerable,
// source: https://github.com/EdOverflow/can-i-take-over-xyz/blob/44e2da47ecb95fc38a0976812fc173e553996189/fingerprints.json
// Cannot test domain below:
//   - agilecrm.com
//   - airee.ru
//   - youtrack.cloud
var TakeoverDomains = []TakeoverDomain{
	{
		ServiceName: "AWS/Elastic Beanstalk",
		Domain:      "us-east-1.elasticbeanstalk.com",
		Type:        NVHO,
	},
	{
		ServiceName: "AWS/S3",
		Domain:      "s3.amazonaws.com",
		Type:        VHO,
		Fingerprint: "NoSuchBucket",
	},
	{
		ServiceName: "Anima",
		Domain:      "animaapp.io",
		Type:        VHO,
		Fingerprint: "Anima - Page Not Found",
	},
	{
		ServiceName: "Bitbucket",
		Domain:      "bitbucket.io",
		Type:        VHO,
		Fingerprint: "Repository not found",
	},
	{
		ServiceName: "Gemfury",
		Domain:      "furyns.com",
		Type:        VHO,
		Fingerprint: "404: This page could not be found.",
	},
	{
		ServiceName: "Ghost",
		Domain:      "ghost.io",
		Type:        VHO,
		Fingerprint: "Domain error",
	},
	{
		ServiceName: "HatenaBlog",
		Domain:      "hatenablog.com",
		Type:        VHO,
		Fingerprint: "404 Blog is not found",
	},
	{
		ServiceName: "Help Juice",
		Domain:      "helpjuice.com",
		Type:        VHO,
		Fingerprint: "We could not find what you're looking for.",
	},
	{
		ServiceName: "Help Scout",
		Domain:      "helpscoutdocs.com",
		Type:        VHO,
		Fingerprint: "No settings were found for this company:",
	},
}
