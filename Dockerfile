FROM alpine:3.6

COPY ./build/linux/altair /usr/local/bin/
COPY ./migration/ /opt/altair/

RUN apk --update upgrade
RUN apk --no-cache add curl tzdata

EXPOSE 2019
ENTRYPOINT ["altair", "run"]
