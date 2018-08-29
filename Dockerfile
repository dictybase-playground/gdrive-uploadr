FROM golang:1.10.4-alpine3.8
LABEL maintainer="Siddhartha Basu<siddhartha-basu@northwestern.edu>"

RUN apk add --no-cache git build-base \
    && go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/dictybase-playground/gdrive-uploadr
WORKDIR /go/src/github.com/dictybase-playground/gdrive-uploadr
COPY Gopkg.* main.go ./
ADD apihelpers apihelpers
ADD auth auth
ADD commands commands
ADD handlers handlers
ADD logger logger
ADD validate validate
RUN dep ensure \
    && go build -o app

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/dictybase-playground/gdrive-uploadr/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
ENV TZ America/Chicago

