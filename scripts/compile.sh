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
    # Generate windows resources
    BUILD_VERSION_COMMA=$(echo -n $BUILD_VERSION_SUFFIX.0 | sed -re 's/^([0-9.]*).*$/\1/' | tr . ,)
    defs=
    [ ! -z $BUILD_VERSION_SUFFIX ] && defs="$defs -D BUILD_VERSION=\"$BUILD_VERSION_SUFFIX\""
	  [ ! -z $BUILD_VERSION_COMMA ]  && defs="$defs -D BUILD_VERSION_COMMA=$BUILD_VERSION_COMMA"

    x86_64-w64-mingw32-windres \
      -i "${__resources}/sloppy.rc" \
      -F pe-x86-64 \
      -o "${__src}/rsrc_amd64.syso" \
      --use-temp-file \
      $defs
    output="${output}.exe"
  fi
fi

echo "Start compiling sloppy ${BUILD_VERSION} for ${GOOS}-${GOARCH}"
go build \
  -ldflags \
    "-X main.GitCommit=${GIT_COMMIT} \
    -X main.Version=${BUILD_VERSION} \
    -X main.VersionPrerelease=${BUILD_PRE_RELEASE}" \
  ${trim_flags:-} \
  -o "${BUILD_TARGET}/${output:-sloppy}" \
  "./src"

# Cleanup
if [ $GOOS = "windows" ]; then
  rm "${__src}/rsrc_amd64.syso"
fi

if [ $? == 0 ]; then
  echo "Output: ${BUILD_TARGET}/${output:-sloppy}"
  echo "==> Successfully"
fi
