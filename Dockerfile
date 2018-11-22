# extending golang to build and test locally and on CI
FROM golang:1.9.3

ENV JQ_VERSION 1.5
ENV JQ_DOWNLOAD_URL https://github.com/stedolan/jq/releases/download/jq-$JQ_VERSION/jq-linux32

RUN curl -fsSL "$JQ_DOWNLOAD_URL" -o jq && \
    chmod +x jq && mv jq /usr/local/bin/jq

# dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Windows Resource generation
RUN go get github.com/josephspurrier/goversioninfo && \
    go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo
