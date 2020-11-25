FROM golang:1.15 as builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-osint/osint
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

RUN mkdir /go/src/osint
WORKDIR /go/src/osint
ADD . .
#ADD ./proto osint/proto

FROM alpine
RUN apk add --no-cache ca-certificates tzdata \
  && cd /tmp/ && wget -q https://github.com/okzk/env-injector/releases/download/v0.0.5/env-injector_0.0.5_linux_amd64.tar.gz \
  && tar xvfz env-injector_0.0.5_linux_amd64.tar.gz && mv env-injector /usr/local/bin/ && rm /tmp/*
COPY --from=builder /go/bin/osint /usr/local/osint/bin/
ENV PORT=18081 \
    DB_MASTER_HOST= \
    DB_MASTER_USER= \
    DB_MASTER_PASSWORD= \
    DB_SLAVE_HOST= \
    DB_SLAVE_USER= \
    DB_SLAVE_PASSWORD= \
    DB_SCHEMA=mimosa \
    DB_PORT=3306 \
    DB_LOG_MODE=false \
    AWS_REGION= \
    AWS_ACCESS_KEY_ID= \
    AWS_SECRET_ACCESS_KEY= \
    AWS_SESSION_TOKEN= \
    SQS_ENDPOINT= \
    PRIVATE_EXPOSE_QUEUE_NAME= \
    PRIVATE_EXPOSE_QUEUE_URL= \
    SUBDOMAIN_QUEUE_NAME= \
    SUBDOMAIN_QUEUE_URL= \
    TZ=Asia/Tokyo
WORKDIR /usr/local/osint
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/osint"]