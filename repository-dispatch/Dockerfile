FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.13.7-alpine3.11 as build
WORKDIR /go/src/github.com/universalbasket/github-actions/repository-dispatch/src
RUN apk update && \
  apk upgrade && \
  apk add --no-cache bash git openssh
RUN adduser -D -g '' user
COPY . /go/src/github.com/universalbasket/github-actions/repository-dispatch/src
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o bin .

FROM scratch
ENV PATH=/bin
ARG DIR
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/universalbasket/github-actions/repository-dispatch/bin /bin/repository-dispatch
COPY --from=build /etc/passwd /etc/passwd
USER user
ENTRYPOINT [ "./bin/repository-dispatch"]
