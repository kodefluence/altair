FROM alpine:3.6

WORKDIR /opt/altair/

COPY ./build_output/linux/altair /usr/local/bin/
COPY ./migration/ /opt/altair/migration/
COPY ./config/ /opt/altair/config/
COPY ./routes/ /opt/altair/routes/
COPY ./env.sample /opt/altair/.env

RUN apk --update upgrade
RUN apk --no-cache add curl tzdata

EXPOSE 1304
ENTRYPOINT ["altair", "run"]
