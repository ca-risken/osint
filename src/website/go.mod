module github.com/ca-risken/osint/src/website

go 1.17

require (
	github.com/aws/aws-sdk-go v1.42.22
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/logging v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/sqs v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20211028073607-f05af5412a86
	github.com/ca-risken/core/proto/finding v0.0.0-20211028073607-f05af5412a86
	github.com/ca-risken/osint/pkg/common v0.0.0-20211112075552-aaad6ef82b13
	github.com/ca-risken/osint/pkg/message v0.0.0-20211112075552-aaad6ef82b13
	github.com/ca-risken/osint/proto/osint v0.0.0-20211112075552-aaad6ef82b13
	github.com/gassara-kys/envconfig v1.4.4
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	google.golang.org/grpc v1.42.0
)

require (
	github.com/Songmu/timeout v0.4.0 // indirect
	github.com/Songmu/wrapcommander v0.1.0 // indirect
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.11.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.24.0 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
