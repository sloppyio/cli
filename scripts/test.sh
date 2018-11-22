#!/usr/bin/env bash
set -e

echo "Start unit tests"

# Running per package for race and coverage option
# Ensures not to test the vendor directory too
packages=$(go list ./... | grep -h -v "/vendor/")
for pkg in ${packages}; do
    go test -race -timeout 30s -covermode=atomic -coverprofile=$(basename ${pkg}).cover ${pkg}
done
# Generate cover profile
echo "mode: atomic" > ./coverage.txt
grep -h -v "^mode:" ./*.cover >> ./coverage.txt && rm ./*.cover

if [ $? == 0 ]; then
  echo "==> Successfull"
fi
