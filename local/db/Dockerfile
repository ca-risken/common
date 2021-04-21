FROM mysql:5.7.32

ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=
ENV MYSQL_ALLOW_EMPTY_PASSWORD=yes

ADD *.sql /docker-entrypoint-initdb.d/
