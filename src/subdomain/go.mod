module github.com/ca-risken/osint/src/subdomain

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.42.22
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/sqs v0.0.0-20211210074045-79fdb4c61950
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20210924100500-e1499111345b
	github.com/ca-risken/core/proto/finding v0.0.0-20211208021136-f6b597adc160
	github.com/ca-risken/osint/pkg/common v0.0.0-20210915063551-7002685890c3
	github.com/ca-risken/osint/pkg/message v0.0.0-20211112065816-37550cc4192d
	github.com/ca-risken/osint/pkg/model v0.0.0-20210915063551-7002685890c3
	github.com/ca-risken/osint/proto/osint v0.0.0-20210915063551-7002685890c3
	github.com/gassara-kys/envconfig v1.4.4
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/miekg/dns v1.1.43
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/valyala/fasthttp v1.30.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20210928044308-7d9f5e0b762b // indirect
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210927142257-433400c27d05 // indirect
	google.golang.org/grpc v1.41.0
)
