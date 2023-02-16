FROM golang:1.19.6-alpine3.16 AS builder
ENV ROOT=/go/src/app
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

ENV GO111MODULE=on
COPY . .

RUN go mod tidy
