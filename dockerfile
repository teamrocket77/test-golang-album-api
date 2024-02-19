FROM mysql

ENV MYSQL_DATABASE=docker
ENV MYSQL_USER=docker
ENV MYSQL_ROOT_PASSWORD=docker
ENV MYSQL_PASSWORD=docker

COPY ./sql/init.sql /docker-entrypoint-initdb.d/
COPY ./sql/main.cnf /etc/my.cnf
