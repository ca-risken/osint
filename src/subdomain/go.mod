module github.com/CyberAgent/mimosa-osint-go/src/subdomain

go 1.15

require (
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20201130105221-b9659eb5f70a
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20201130105221-b9659eb5f70a
	github.com/CyberAgent/mimosa-osint/pkg/common v0.0.0-20201211055630-9e7fa6b8b56a
	github.com/CyberAgent/mimosa-osint/pkg/message v0.0.0-20201211055630-9e7fa6b8b56a
	github.com/CyberAgent/mimosa-osint/pkg/model v0.0.0-20201211055630-9e7fa6b8b56a
	github.com/CyberAgent/mimosa-osint/proto/osint v0.0.0-20201211055630-9e7fa6b8b56a
	github.com/aws/aws-sdk-go v1.37.10
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f
	github.com/go-sql-driver/mysql v1.5.0
	github.com/h2ik/go-sqs-poller/v3 v3.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/miekg/dns v1.1.35
	github.com/sirupsen/logrus v1.7.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sys v0.0.0-20201211002650-1f0c578a6b29 // indirect
	google.golang.org/genproto v0.0.0-20201210142538-e3217bee35cc // indirect
	google.golang.org/grpc v1.34.0
)
