FROM golang:1.10.3-alpine

MAINTAINER Ryan Holcombe <rholcombe30@gmail.com>

ARG VCS_REF
ARG BUILD_DATE

# Metadata
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/sythe21/s3api" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.docker.dockerfile="/Dockerfile"

COPY . /go/src/github.com/sythe21/s3api

ENV GIT_SHA $VCS_REF
ENV GOPATH /go
RUN cd $GOPATH/src/github.com/sythe21/s3api && go install -v .

CMD ["s3api"]

EXPOSE 8888
	
