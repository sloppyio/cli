FROM golang:1.9.3

ENV JQ_VERSION 1.5
ENV JQ_DOWNLOAD_URL https://github.com/stedolan/jq/releases/download/jq-$JQ_VERSION/jq-linux32

RUN apt-get -y -q update && \
    apt-get -y -q install httpie && \
    curl -fsSL "$JQ_DOWNLOAD_URL" -o jq && \
    chmod +x jq && mv jq /usr/local/bin/jq && \
    # Windows Resource generation
    go get github.com/josephspurrier/goversioninfo && \
    go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo && \
    # Bats
    curl -sSL https://github.com/sstephenson/bats/archive/v0.4.0.tar.gz -o bats.tar.gz && \
    tar -xf bats.tar.gz && ./bats-0.4.0/install.sh /usr/local && rm -rf bats-0.4.0
