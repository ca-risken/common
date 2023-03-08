FROM public.ecr.aws/risken/base/risken-base:v0.0.1 as base

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
COPY --from=base /usr/local/bin/env-injector /usr/local/bin/
