module github.com/ca-risken/osint/src/osint

go 1.16

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aws/aws-sdk-go v1.42.3
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/database v0.0.0-20211028082529-ca6c2d1ed0ee
	github.com/ca-risken/common/pkg/sqs v0.0.0-20211028082529-ca6c2d1ed0ee // indirect
	github.com/ca-risken/common/pkg/xray v0.0.0-20211028082529-ca6c2d1ed0ee
	github.com/ca-risken/osint/pkg/message v0.0.0-20211108025335-540430e56cd1
	github.com/ca-risken/osint/pkg/model v0.0.0-20211108025335-540430e56cd1
	github.com/ca-risken/osint/proto/osint v0.0.0-20211108025335-540430e56cd1
	github.com/gassara-kys/envconfig v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20211111160137-58aab5ef257a // indirect
	golang.org/x/sys v0.0.0-20211111213525-f221eed1c01e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211111162719-482062a4217b // indirect
	google.golang.org/grpc v1.42.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gorm.io/driver/mysql v1.1.3 // indirect
	gorm.io/gorm v1.22.2
)
