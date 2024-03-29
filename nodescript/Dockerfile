FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.19-alpine3.17 as build
WORKDIR /go/src/github.com/ubio/github-actions/nodescript
RUN apk update && \
  apk upgrade && \
  apk add --no-cache bash git openssh
RUN adduser -D -g '' user
COPY ./src /go/src/github.com/ubio/github-actions/nodescript
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o bin .

FROM scratch
ENV PATH=/bin
ARG DIR
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/ubio/github-actions/nodescript/bin /bin/nodescript
COPY --from=build /etc/passwd /etc/passwd
USER user
ENTRYPOINT [ "/bin/nodescript"]
