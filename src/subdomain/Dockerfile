FROM golang:1.15 AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-osint-go/src/subdomain/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM alpine
COPY theHarvester /theHarvester
RUN apk add --no-cache ca-certificates tzdata build-base python3 py-pip python3-dev libffi-dev gcc libc-dev libxml2-dev libxslt-dev \
  && pip install wheel \
  && pip install -r /theHarvester/requirements/base.txt \
  && cd /tmp/ && wget -q https://github.com/okzk/env-injector/releases/download/v0.0.5/env-injector_0.0.5_linux_amd64.tar.gz \
  && tar xvfz env-injector_0.0.5_linux_amd64.tar.gz && mv env-injector /usr/local/bin/ && rm /tmp/* \
  && mkdir /results
COPY --from=builder /go/bin/subdomain /usr/local/subdomain/bin/
ENV PORT=19002 \
    AWS_REGION= \
    AWS_ACCESS_KEY_ID= \
    AWS_SECRET_ACCESS_KEY= \
    AWS_SESSION_TOKEN= \
    SQS_ENDPOINT= \
    SUBDOMAIN_SEARCH_QUEUE_URL= \
    MAX_NUMBER_OF_MESSAGE=10 \
    WAIT_TIME_SECOND=20 \
    FINDING_SVC_ADDR= \
    ALERT_SVC_ADDR= \
    OSINT_SVC_ADDR= \
    RESULT_PATH=/results \
    HARVESTER_PATH=/theHarvester \
    TZ=Asia/Tokyo

WORKDIR /usr/local/subdomain
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/subdomain"]
