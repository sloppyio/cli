#! /bin/bash
set -e

echo "Start fetching dependencies"

go get github.com/golang/lint/golint
go get gopkg.in/yaml.v2
go get github.com/olekukonko/tablewriter
go get github.com/mitchellh/cli
go get github.com/hashicorp/go-version

if [ $? == 0 ]; then
  echo "==> Successfully"
fi
