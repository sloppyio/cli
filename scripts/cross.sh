#!/usr/bin/env bash
set -e

source "${__scripts}/dependencies.sh"

echo "Start cross-compiling sloppy"
echo
for GOOS in darwin linux windows; do
  export GOOS=$GOOS
  for GOARCH in amd64; do
    export GOARCH=$GOARCH
    source "${__scripts}/compile.sh"
  done
done

if [ $? == 0 ]; then
  echo
  echo "Output: ${BUILD_TARGET}"
  echo "==> Successfully"
fi
