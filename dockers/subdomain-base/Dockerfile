FROM golang:1.18.2 AS builder
ENV CGO_ENABLED=0
ARG ENV_INJECTOR_VERSION=v0.0.6

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM ubuntu:20.04
ARG THE_HARVESSTER_VERSION=4.0.0
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -qy python3 python3-pip libffi-dev libxml2-dev libxslt1-dev git wget  \ 
  && /usr/bin/python3 -m pip install --upgrade pip && apt clean && apt autoremove -qy
RUN git clone https://github.com/laramies/theHarvester.git -b ${THE_HARVESSTER_VERSION} \
  && pip install --no-cache-dir -r /theHarvester/requirements/base.txt \
  && mkdir /results
COPY --from=builder /go/bin/env-injector /usr/local/bin/
