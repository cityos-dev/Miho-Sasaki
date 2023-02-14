FROM golang:1.17-alpine AS builder
ENV ROOT=/go/src/app
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

ADD . /myapp
WORKDIR /myapp

RUN CGO_ENABLED=0 GOOS=linux go build -o server ../myapp/main.go

FROM alpine:3.10
COPY --from=builder /myapp/server /app
#EXPOSE 8080
#CMD ["/app"]
