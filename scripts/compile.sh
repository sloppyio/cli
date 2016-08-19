#!/usr/bin/env bash
set -e

# Set GOOS and GOARCH to HOST
if [ -z "$GOOS" ]; then
  GOOS=`go env GOOS`
fi
if [ -z "$GOARCH" ]; then
  GOARCH=`go env GOARCH`
fi

case "${1:-}" in
"test")
  ;;
"beta")
  ;;
"rc.[0-9]*")
  ;;
"release")
  # Remove actual path from source files.
  trim_flags="-gcflags=-trimpath=${GOPATH} \
              -asmflags=-trimpath=${GOPATH}"
  GOROOT_FINAL=/usr/local/go
  ;;
*)
  set -- $1 "dev"
esac
source "${__scripts}/.version.sh"

# Name output file
if [ "${output}" != "sloppy" ]; then
  output="sloppy-${GOOS}-${GOARCH}"
  if [ $GOOS = "windows" ]; then
    output="${output}.exe"
  fi
fi

echo "Start compiling sloppy ${BUILD_VERSION} for ${GOOS}-${GOARCH}"
go build \
  -ldflags \
    "-X main.GitCommit=${GIT_COMMIT} \
    -X main.Version=${BUILD_VERSION} \
    -X main.VersionPrerelease=${pre_release_version}" \
  ${trim_flags:-} \
  -o "${BUILD_TARGET}/${output:-sloppy}" \
  "./src"

if [ $? == 0 ]; then
  echo "Output: ${BUILD_TARGET}/${output:-sloppy}"
  echo "==> Successfully"
fi
