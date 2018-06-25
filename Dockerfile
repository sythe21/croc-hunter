FROM golang:1.10.3-alpine

ARG GIT_COMMIT
ARG VERSION
ARG BUILD_DATE
LABEL REPO="https://github.com/sythe21/s3api"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION
LABEL BUILD_DATE=$BUILD_DATE

ADD . /go/src/github.com/sythe21/s3api
WORKDIR /go/src/github.com/sythe21/s3api

RUN apk update && apk add make git
RUN make build

ENTRYPOINT ["/go/bin/s3api"]