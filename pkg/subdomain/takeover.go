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

func checkTakeover(h host) takeover {
	cname := resolveCName(h.HostName)
	if cname == "" {
		return takeover{}
	}
	return takeover{
		Domain: h.HostName,
		CName:  cname,
		IsDown: h.isDown(),
	}
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

func (t *takeover) makeFinding(projectID uint32, dataSource string) (*finding.FindingForUpsert, error) {
	if zero.IsZeroVal(*t) {
		return nil, nil
	}
	data, err := json.Marshal(map[string]takeover{"data": *t})
	if err != nil {
		return nil, err
	}
	finding := &finding.FindingForUpsert{
		Description:      t.getDescription(),
		DataSource:       dataSource,
		DataSourceId:     generateDataSourceID(fmt.Sprintf("%v_%v", t.Domain, t.CName)),
		ResourceName:     t.Domain,
		ProjectId:        projectID,
		OriginalScore:    t.getScore(),
		OriginalMaxScore: 10.0,
		Data:             string(data),
	}
	return finding, nil
}

func (t *takeover) getScore() float32 {
	var score float32
	if t.IsDown {
		score = 6.0
		if t.matchTakeoverList() {
			score = score + 2.0
		}
		return score
	}
	return 1.0
}

func (t *takeover) getDescription() string {
	var desc string
	if t.IsDown {
		desc = fmt.Sprintf("%s seems to be down. It has subdomain takeover risk.", t.Domain)
		desc = desc + fmt.Sprintf("(CName: %s)", t.CName)
	} else {
		desc = fmt.Sprintf("%s has a CName record.", t.Domain)
		desc = desc + fmt.Sprintf("(CName: %s)", t.CName)
	}
	if len(desc) > 200 {
		desc = desc[:196] + " ..." // cut long text
	}
	return desc
}

func (t *takeover) matchTakeoverList() bool {
	takeoverList := GetTakeOverList()
	for _, takeoverDomain := range takeoverList {
		if strings.Contains(t.CName, takeoverDomain) {
			return true
		}
	}
	return false
}

type takeover struct {
	Domain string `json:"domain"`
	CName  string `json:"forwarding_domain"`
	IsDown bool   `json:"is_down"`
}
