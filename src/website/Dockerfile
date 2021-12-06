FROM golang:1.16 as builder
WORKDIR /go/src/github.com/ca-risken/osint/src/website/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GOPRIVATE="github.com/CyberAgent/*,github.com/ca-risken/*"
ARG ENV_INJECTOR_VERSION=v0.0.6
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v
RUN go build && cp website /go/bin/

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM node:14-alpine
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true
ENV CHROMIUM_BIN /usr/bin/chromium-browser
RUN  apk update && apk add -u --no-cache \
  git \
	nodejs \
  udev \
  chromium \
  ttf-freefont \
  yarn \
  tini \
  && cd /opt \
  && git clone https://github.com/AliasIO/wappalyzer.git -b v6.9.7 && cd wappalyzer \
  && yarn install \
  && yarn run link
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/website /usr/local/website/bin/
ENV DEBUG= \
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
  AWS_XRAY_SDK_DISABLED=TRUE \
  AWS_XRAY_CONTEXT_MISSING=LOG_ERROR \
  RESULT_DIR=/tmp \
  TZ=Asia/Tokyo
WORKDIR /usr/local/website
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/env-injector"]
CMD ["bin/website"]

