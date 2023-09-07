FROM mysql:8.0.34

ENV MYSQL_ROOT_PASSWORD=
ENV MYSQL_ALLOW_EMPTY_PASSWORD=yes

COPY 0*.sql /docker-entrypoint-initdb.d/
