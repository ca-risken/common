FROM public.ecr.aws/risken/base/risken-base:v0.0.1 as base

FROM ubuntu:20.04
ARG THE_HARVESSTER_VERSION=4.0.0
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -qy python3 python3-pip libffi-dev libxml2-dev libxslt1-dev git wget  \ 
  && /usr/bin/python3 -m pip install --upgrade pip && apt clean && apt autoremove -qy
RUN git clone https://github.com/laramies/theHarvester.git -b ${THE_HARVESSTER_VERSION} \
  && pip install --no-cache-dir -r /theHarvester/requirements/base.txt \
  && mkdir /results
COPY --from=base /usr/local/bin/env-injector /usr/local/bin/
