FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-osint/src/subdomain/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM ubuntu:hirsute
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -qy gnupg2 \
  && apt-get update && apt-get install -qy python3 python3-pip libffi-dev libxml2-dev libxslt1-dev git wget  \ 
  && /usr/bin/python3 -m pip install --upgrade pip && apt-get clean && apt-get autoremove -qy
RUN git clone https://github.com/laramies/theHarvester.git -b 4.0.0 \
  && pip install --no-cache-dir -r /theHarvester/requirements/base.txt
RUN cd /tmp/ && wget -q https://github.com/gassara-kys/env-injector/releases/download/v0.0.6/env-injector_0.0.6_Linux_x86_64.tar.gz \
  && tar xvfz env-injector_0.0.6_Linux_x86_64.tar.gz && mv env-injector /usr/local/bin/ && rm -rf /tmp/* \
  && mkdir /results
COPY --from=builder /go/bin/subdomain /usr/local/subdomain/bin/
ENV PORT=19003 \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  SQS_ENDPOINT= \
  MAX_NUMBER_OF_MESSAGE=5 \
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
