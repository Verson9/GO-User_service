FROM golang:1.17.7-alpine3.15
MAINTAINER Przemysław Surma <przemyslaw.surma@outlook.com>

ENV DBPASS=password
RUN mkdir $GOPATH/app
WORKDIR $GOPATH/app
COPY ../go.mod .
COPY ../go.sum .
COPY user-service .
RUN go mod vendor
CMD go build -o user-service/cmd/user-service/main.go
