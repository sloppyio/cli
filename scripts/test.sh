#!/usr/bin/env bash
set -e

echo "Start unit tests"

# Running per package for race and coverage option
# Ensures not to test the vendor directory too
packages=$(go list ./... | grep -h -v "/vendor/")
for pkg in ${packages}; do
    go test -v -race -timeout 30s -covermode=count -coverprofile=$(basename ${pkg}).cover ${pkg}
done
# Generate cover profile
echo "mode: count" > ./coverage.txt
grep -h -v "^mode:" ./*.cover >> ./coverage.txt && rm ./*.cover

if [ $? == 0 ]; then
  echo "==> Successfully"
fi

echo "Start integration tests"

set -- $1 "test"
output="sloppy"
source "${__scripts}/compile.sh"

export SLOPPY_API_URL=https://api.sloppy.io/v1/
export PATH="${BUILD_TARGET}:$PATH"

{
  bats --version > /dev/null
} || {
  echo "Integration tests are not supported by your OS."
  exit 1
}

bats ${__tests}/integration/cli/tests.bats
bats ${__tests}/integration/cli/volume.bats

unset SLOPPY_APITOKEN
bats ${__tests}/integration/cli/nologin.bats


if [ $? == 0 ]; then
  echo "Output: ${BUILD_TARGET}/${output:-sloppy}"
  echo "==> Successfully"
fi
