FROM mysql:5.7

ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=
ENV MYSQL_ALLOW_EMPTY_PASSWORD=yes

COPY 0*.sql /docker-entrypoint-initdb.d/
