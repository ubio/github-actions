FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.14.2-alpine3.11 as build
WORKDIR /go/src/github.com/universalbasket/github-actions/pull-request
RUN apk update && \
  apk upgrade && \
  apk add --no-cache bash git openssh
RUN adduser -D -g '' user
COPY ./src /go/src/github.com/universalbasket/github-actions/pull-request
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o bin .

FROM scratch
ENV PATH=/bin
ARG DIR
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/universalbasket/github-actions/pull-request/bin /bin/pull-request
COPY --from=build /etc/passwd /etc/passwd
USER user
ENTRYPOINT [ "/bin/pull-request"]
