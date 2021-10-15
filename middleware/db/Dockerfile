# MySQL5.7 is not supported for `arm64` arch.
FROM mysql:latest

ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=
ENV MYSQL_ALLOW_EMPTY_PASSWORD=yes

COPY 0*.sql /docker-entrypoint-initdb.d/
