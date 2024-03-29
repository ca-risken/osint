FROM golang:1.18.2 as builder
WORKDIR /go/src/github.com/ca-risken/osint
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/website cmd/website/main.go

FROM public.ecr.aws/risken/osint/website-base:v0.0.2
COPY --from=builder /go/bin/website /usr/local/website/bin/
ENV DEBUG= \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  SQS_ENDPOINT= \
  WEBSITE_QUEUE_NAME= \
  WEBSITE_QUEUE_URL= \
  MAX_NUMBER_OF_MESSAGE= \
  WAIT_TIME_SECOND=20 \
  FINDING_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  AWS_SVC_ADDR= \
  RESULT_DIR=/tmp \
  TZ=Asia/Tokyo
WORKDIR /usr/local/website
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/env-injector"]
CMD ["bin/website"]

