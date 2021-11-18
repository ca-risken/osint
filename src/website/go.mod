module github.com/ca-risken/osint/src/website

go 1.16

require (
	github.com/aws/aws-sdk-go v1.42.3
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/sqs v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20211028073607-f05af5412a86
	github.com/ca-risken/core/proto/finding v0.0.0-20211028073607-f05af5412a86
	github.com/ca-risken/osint/pkg/common v0.0.0-20211112075552-aaad6ef82b13
	github.com/ca-risken/osint/pkg/message v0.0.0-20211112075552-aaad6ef82b13
	github.com/ca-risken/osint/proto/osint v0.0.0-20211112075552-aaad6ef82b13
	github.com/gassara-kys/envconfig v1.4.4
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/sirupsen/logrus v1.8.1
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	google.golang.org/grpc v1.42.0
)
