FROM golang:alpine

ENV GOPATH /app

RUN apk update
RUN apk add git
RUN apk add --no-cache openssl
RUN go get -u github.com/kardianos/govendor
RUN go install -v github.com/kardianos/govendor

RUN pwd
RUN ls

COPY ./common /app/src/saveIdent.com/common/
COPY ./vendor/vendor.json /app/src/saveIdent.com/vendor/vendor.json

COPY ./server/activeMonitorService /app/src/saveIdent.com/server/activeMonitorService
COPY ./server/deviceInputService /app/src/saveIdent.com/server/deviceInputService
COPY ./server/sqlInterfaceService /app/src/saveIdent.com/server/sqlInterfaceService

WORKDIR /app/src/saveIdent.com/
RUN /app/bin/govendor sync

WORKDIR /app/src/saveIdent.com/common/elasticService
RUN /app/bin/govendor sync

WORKDIR /app

RUN go install -v saveIdent.com/server/activeMonitorService

CMD bin/deviceInputService

EXPOSE 1992