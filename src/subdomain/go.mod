module github.com/ca-risken/osint/src/subdomain

go 1.16

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.42.2
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20211028082529-ca6c2d1ed0ee
	github.com/ca-risken/common/pkg/sqs v0.0.0-20211028082529-ca6c2d1ed0ee
	github.com/ca-risken/common/pkg/xray v0.0.0-20211028082529-ca6c2d1ed0ee
	github.com/ca-risken/core/proto/alert v0.0.0-20211028073607-f05af5412a86
	github.com/ca-risken/core/proto/finding v0.0.0-20211028073607-f05af5412a86
	github.com/ca-risken/osint/pkg/common v0.0.0-20211108025335-540430e56cd1
	github.com/ca-risken/osint/pkg/message v0.0.0-20211108025335-540430e56cd1
	github.com/ca-risken/osint/pkg/model v0.0.0-20211108025335-540430e56cd1
	github.com/ca-risken/osint/proto/osint v0.0.0-20211108025335-540430e56cd1
	github.com/gassara-kys/envconfig v1.4.4
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/miekg/dns v1.1.43
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20211111083644-e5c967477495 // indirect
	golang.org/x/sys v0.0.0-20211110154304-99a53858aa08 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247 // indirect
	google.golang.org/grpc v1.42.0
)
