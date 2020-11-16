.PHONY: all install clean network fmt build doc
all: run

install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/go-grpc-middleware

clean:
	rm -f pkg/pb/**/*.pb.go
	rm -f doc/*.md

# @see https://github.com/CyberAgent/mimosa-common/tree/master/local
network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

fmt: proto/**/*.proto
	clang-format -i proto/**/*.proto

doc: fmt
	protoc \
		--proto_path=proto \
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

build: fmt
	protoc \
		--proto_path=proto \
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--go_out=plugins=grpc,paths=source_relative:proto proto/**/*.proto \
		proto/**/*.proto;

go-test: build
	cd proto/osint  && go test ./...
	cd pkg/message      && go test ./...
	cd src/osint    && go test ./...
	cd src/private-expose         && go test ./...

go-mod-update:
	cd src/osint \
		&& go get -u \
			github.com/CyberAgent/mimosa-osint/...
	cd src/private-expose \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/... \
			github.com/CyberAgent/mimosa-osint/...

go-mod-tidy: build
	cd pkg/common   && go mod tidy
	cd pkg/model   && go mod tidy
	cd pkg/message   && go mod tidy
	cd src/osint && go mod tidy
	cd src/private-expose      && go mod tidy

run: go-test network
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down
