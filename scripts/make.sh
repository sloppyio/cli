#!/usr/bin/env bash
set -e

# Magic folders
__dir="$(cd "$(dirname "{BASH_SOURCE[0]}")" && pwd)"
__src="${__dir}"
__scripts="${__dir}/scripts"
__tests="${__dir}/tests"
__bundles="${__dir}/bundles"
__resources="${__dir}/.resources-windows"

arg=${1:-}
shift

case $arg in
"build")
  source "${__scripts}/compile.sh"
  ;;
"cross")
  source "${__scripts}/cross.sh"
  ;;
"deploy")
  source "${__scripts}/cross.sh"
  echo "Copy to ${__bundles}/latest"
  mkdir -p ${__bundles}/latest
  cp -a ${__bundles}/${BUILD_VERSION_SUFFIX}/. ${__bundles}/latest
  ;;
"test")
  source "${__scripts}/test.sh"
  ;;
*)
  echo "${arg} is not supported"
esac
