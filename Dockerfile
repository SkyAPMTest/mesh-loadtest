FROM golang:latest
LABEL maintainer="hanahmily@apache.org"
WORKDIR $GOPATH/src/github.com/SkyAPMTest/mesh-loadtest
ADD . $GOPATH/src/github.com/SkyAPMTest/mesh-loadtest
RUN go build .

ENTRYPOINT  ["./mesh-loadtest"]
