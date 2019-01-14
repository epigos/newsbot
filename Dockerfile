FROM golang:1.10.2

WORKDIR /go/src/newsbot

ADD . /go/src/github.com/epigos/newsbot/

RUN make build

EXPOSE 5050
EXPOSE 443

ENTRYPOINT /go/bin/newsbot