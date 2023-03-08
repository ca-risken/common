FROM public.ecr.aws/risken/base/risken-base:v0.0.1 as base

FROM aquasec/trivy:0.31.2
RUN apk add --no-cache ca-certificates tzdata
COPY --from=base /usr/local/bin/env-injector /usr/local/bin/
