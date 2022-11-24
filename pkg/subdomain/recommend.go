package subdomain

type recommend struct {
	Type           string `json:"type,omitempty"`
	Risk           string `json:"risk,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

func getRecommend(category string) recommend {
	return recommendMap[category]
}

var recommendMap = map[string]recommend{
	"Takeover": {
		Type: "Domain/TakeOver",
		Risk: `Domain Takeover Lisk
		- The target domain has a CNAME.
		- If it is down or not under your control, it may be hijacked.`,
		Recommendation: `Delete unused DNS records.
		- https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/02-Configuration_and_Deployment_Management_Testing/10-Test_for_Subdomain_Takeover`,
	},
	"PrivateExpose": {
		Type: "Domain/PrivateExpose",
		Risk: `Site is Exposed to Private
		- The target site is open to the public.
		- If a site that does not need to be public is public, it may be compromised or hijacked.`,
		Recommendation: `Please restrict access by IP address or authentication.`,
	},
	"CertificateExpiration": {
		Type: "Domain/CertificateExpiration",
		Risk: `Certificate 
		- Expired certificates may allow for spoofing attacks.
		- Corporate brands are adversely affected and business is put at risk.
		- Many browsers will display a warning due to certificate expiration`,
		Recommendation: `Update your certificates before the expiration date.`,
	},
}
