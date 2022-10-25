FROM golang:1.18.2 as builder
ENV CGO_ENABLED=0
ARG ENV_INJECTOR_VERSION=v0.0.6

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM node:14-alpine
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true
ENV CHROMIUM_BIN /usr/bin/chromium-browser
ARG WAPPALYZER_VERSION=v6.9.7
RUN  apk update && apk add -u --no-cache \
  git \
  nodejs \
  udev \
  chromium \
  ttf-freefont \
  yarn \
  tini \
  && cd /opt \
  && git clone https://github.com/AliasIO/wappalyzer.git -b ${WAPPALYZER_VERSION} \
  && cd wappalyzer \
  && yarn install \
  && yarn run link
COPY --from=builder /go/bin/env-injector /usr/local/bin/
