FROM golang:latest
LABEL maintainer="hanahmily@apache.org"
WORKDIR $GOPATH/src/github.com/SkyWalkingTest/mesh-loadtest
ADD . $GOPATH/src/github.com/SkyWalkingTest/mesh-loadtest
RUN go build .

ENTRYPOINT  ["./mesh-loadtest"]
