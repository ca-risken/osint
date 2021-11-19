module github.com/ca-risken/osint/src/osint

go 1.16

require (
	github.com/aws/aws-sdk-go v1.40.48
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/database v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/rpc v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/osint/pkg/message v0.0.0-20211112065816-37550cc4192d
	github.com/ca-risken/osint/pkg/model v0.0.0-20210908024505-bad8297bda4e
	github.com/ca-risken/osint/proto/osint v0.0.0-20210908024505-bad8297bda4e
	github.com/gassara-kys/envconfig v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210708141623-e76da96a951f // indirect
	google.golang.org/grpc v1.42.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gorm.io/gorm v1.21.12
)
