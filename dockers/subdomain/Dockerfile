FROM golang:1.18.2 AS builder
WORKDIR /go/src/github.com/ca-risken/osint/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/subdomain cmd/subdomain/main.go

FROM public.ecr.aws/risken/osint/subdomain-base:v0.0.1
COPY --from=builder /go/bin/subdomain /usr/local/subdomain/bin/
ENV PORT=19003 \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  SQS_ENDPOINT= \
  MAX_NUMBER_OF_MESSAGE= \
  WAIT_TIME_SECOND=20 \
  FINDING_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  OSINT_SVC_ADDR= \
  RESULT_PATH=/results \
  HARVESTER_PATH=/theHarvester \
  SUBDOMAIN_QUEUE_NAME= \
  SUBDOMAIN_QUEUE_URL= \
  TZ=Asia/Tokyo

WORKDIR /usr/local/subdomain
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/subdomain"]
