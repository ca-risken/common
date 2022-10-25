FROM golang:1.18.2 as builder
ENV CGO_ENABLED=0
ARG ENV_INJECTOR_VERSION=v0.0.6

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM node:lts-alpine3.12
ARG CLOUDSPLOIT_VERSION=v2.0.0
RUN apk add --no-cache ca-certificates tzdata git \
  && cd /opt  \
  && git clone https://github.com/aquasecurity/cloudsploit.git -b ${CLOUDSPLOIT_VERSION} --depth 1 \
  && cd cloudsploit \
  && yarn install  \
  && chmod +x index.js
COPY --from=builder /go/bin/env-injector /usr/local/bin/
