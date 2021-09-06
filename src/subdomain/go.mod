module github.com/CyberAgent/mimosa-osint-go/src/subdomain

go 1.16

require (
	github.com/CyberAgent/mimosa-osint/pkg/common v0.0.0-20210709011430-ef6f10e6a89a
	github.com/CyberAgent/mimosa-osint/pkg/message v0.0.0-20210709011430-ef6f10e6a89a
	github.com/CyberAgent/mimosa-osint/pkg/model v0.0.0-20210709011430-ef6f10e6a89a
	github.com/CyberAgent/mimosa-osint/proto/osint v0.0.0-20210709011430-ef6f10e6a89a
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.39.3
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20210906122657-d2be54cc7181
	github.com/ca-risken/common/pkg/xray v0.0.0-20210906122657-d2be54cc7181
	github.com/ca-risken/core/proto/alert v0.0.0-20210906115102-3cabd5f9511a
	github.com/ca-risken/core/proto/finding v0.0.0-20210906115102-3cabd5f9511a
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/miekg/dns v1.1.42
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210708141623-e76da96a951f // indirect
	google.golang.org/grpc v1.39.0
)
