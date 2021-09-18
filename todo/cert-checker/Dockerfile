FROM golang:1.14 as build
RUN go get github.com/genkiroid/cert
WORKDIR /go/src/github.com/genkiroid/cert
RUN go get -d -v ./...
RUN go install -v ./...
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh"]
