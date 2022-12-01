FROM alpine:3.17

WORKDIR /opt/altair/

COPY ./build_output/linux/altair /usr/local/bin/
COPY ./config/ /opt/altair/config/
COPY ./routes/ /opt/altair/routes/
COPY ./env.sample /opt/altair/.env

RUN apk --update upgrade
RUN apk --no-cache add curl tzdata

EXPOSE 1304
ENTRYPOINT ["altair", "run"]
