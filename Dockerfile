FROM golang:1.17.7
LABEL MAINTAINER Przemysław Surma <przemyslaw.surma@outlook.com>
RUN mkdir $GOPATH/GO-User_service
WORKDIR $GOPATH/GO-User_service

COPY go.mod .
COPY go.sum .
COPY user-service user-service

RUN go mod vendor
RUN go build -o app ./user-service/cmd/user-service/main.go
ENTRYPOINT ["./app"]
