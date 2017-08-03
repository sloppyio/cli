#!/usr/bin/env bash
set -e

echo "Start unit tests"
go test -v ./...
if [ $? == 0 ]; then
  echo "==> Successfully"
fi

echo "Start integration tests"

set -- $1 "test"
output="sloppy"
source "${__scripts}/compile.sh"

export SLOPPY_APIHOST=https://api.sloppy.io
export PATH="${BUILD_TARGET}:$PATH"

{
  ${__tests}/integration/cli/bats-0.4.0/libexec/bats --version > /dev/null
} || {
  echo "Integration tests are not supported by your OS."
  exit 1
}

${__tests}/integration/cli/bats-0.4.0/libexec/bats ${__tests}/integration/cli/tests.bats
${__tests}/integration/cli/bats-0.4.0/libexec/bats ${__tests}/integration/cli/volume.bats

unset SLOPPY_APITOKEN
${__tests}/integration/cli/bats-0.4.0/libexec/bats ${__tests}/integration/cli/nologin.bats


if [ $? == 0 ]; then
  echo "Output: ${BUILD_TARGET}/${output:-sloppy}"
  echo "==> Successfully"
fi
