package subdomain

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/ca-risken/core/proto/finding"
	"github.com/miekg/dns"
	"github.com/vikyd/zero"
)

func searchTakeover(domain string) takeover {
	cname := resolveCName(domain)
	if !zero.IsZeroVal(cname) {
		return takeover{Domain: domain, CName: cname}
	}
	return takeover{}
}

func resolveCName(domain string) string {
	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeCNAME)
	m.RecursionDesired = true
	// Only normally accessible domains, exclude temporarily inaccessible ex. service unavailable, are scanned, so error is ignored.
	r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53"))
	if err != nil {
		return ""
	}
	if zero.IsZeroVal(r.Answer) {
		return ""
	}
	return r.Answer[0].(*dns.CNAME).Target
}

func (c *takeover) makeFinding(isDown bool, projectID uint32, dataSource string) (*finding.FindingForUpsert, error) {
	if zero.IsZeroVal(*c) {
		return nil, nil
	}
	score := c.getScore(isDown)
	description := c.getDescription(isDown)
	data, err := json.Marshal(map[string]takeover{"data": *c})
	if err != nil {
		return nil, err
	}
	finding := &finding.FindingForUpsert{
		Description:      description,
		DataSource:       dataSource,
		DataSourceId:     generateDataSourceID(fmt.Sprintf("%v_%v", c.Domain, c.CName)),
		ResourceName:     c.Domain,
		ProjectId:        projectID,
		OriginalScore:    score,
		OriginalMaxScore: 10.0,
		Data:             string(data),
	}
	return finding, nil
}

func (c *takeover) getScore(isDown bool) float32 {
	var score float32
	if isDown {
		score = 6.0
		if c.matchTakeoverList() {
			score = score + 2.0
		}
		return score
	}
	return 1.0
}

func (c *takeover) getDescription(isDown bool) string {
	var desc string
	if isDown {
		desc = fmt.Sprintf("%s seems to be down. It has subdomain takeover risk.", c.Domain)
		desc = desc + fmt.Sprintf("(CName: %s)", c.CName)
	} else {
		desc = fmt.Sprintf("%s has a CName record.", c.Domain)
		desc = desc + fmt.Sprintf("(CName: %s)", c.CName)
	}
	if len(desc) > 200 {
		desc = desc[:196] + " ..." // cut long text
	}
	return desc
}

func (c *takeover) matchTakeoverList() bool {
	takeoverList := GetTakeOverList()
	for _, takeoverDomain := range takeoverList {
		if strings.Contains(c.CName, takeoverDomain) {
			return true
		}
	}
	return false
}

type takeover struct {
	Domain string `json:"domain"`
	CName  string `json:"forwarding_domain"`
}
