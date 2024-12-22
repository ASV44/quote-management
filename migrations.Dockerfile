FROM alpine:latest
RUN apk update && apk add --no-cache curl

USER root

WORKDIR /

# Install migrate for running DB migrations
RUN curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate

# Add wait-for script for local
ENV WAIT_VERSION v2.2.3
ADD https://github.com/eficode/wait-for/releases/download/$WAIT_VERSION/wait-for /wait-for
RUN chmod +x /wait-for

## Copy migrations
COPY db/migrations /migrations

ADD db/db_entrypoint.sh /
RUN chmod +x /db_entrypoint.sh

ENTRYPOINT [ "sh", "/db_entrypoint.sh" ]

CMD ["up"]