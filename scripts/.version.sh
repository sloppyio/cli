#!/usr/bin/env bash
BUILD_VERSION=$(< ./version.txt)
BUILD_VERSION_SUFFIX=$BUILD_VERSION

# Pre release version
if [ "${1:-}" != "release" ]; then
  BUILD_VERSION_SUFFIX="${BUILD_VERSION}-${1:-}"
  BUILD_PRE_RELEASE="${1:-}"
fi

BUILD_TARGET="${__bundles}/$BUILD_VERSION_SUFFIX"

if [ ! -d "$BUILD_TARGET" ]; then
  mkdir -p $BUILD_TARGET
fi

if [ -n ${GIT_COMMIT:-} ]; then
  GIT_COMMIT=$(git rev-parse HEAD || "")
fi
