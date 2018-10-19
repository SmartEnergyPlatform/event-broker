FROM golang:1.11


COPY . /go/src/event-broker
WORKDIR /go/src/event-broker

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./event-broker