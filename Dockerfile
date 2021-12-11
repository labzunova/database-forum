FROM golang:1.13 AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main .

EXPOSE 5000

ENV PGPASSWORD 1111

#CMD service postgresql start && psql -h localhost -d forums -U lbznv -p 5432 -a -q -f ./init.sql && ./main
