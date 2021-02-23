#!/usr/bin/env bash

set -euo pipefail

log() {
  local date
  local level
  local debug
  local line
  local msg
  local div
  local -A colours

  level="${1}"
  msg="${*:2}"

  date="$(date '+%T')"
  debug="${DEBUG:-0}"

  colours['DEBUG']='\033[34m'  # Blue
  colours['INFO']='\033[32m'   # Green
  colours['WARN']='\033[33m'   # Yellow
  colours['ERROR']='\033[31m'  # Red
  colours['DEFAULT']='\033[0m' # Uncoloured

  div="================================================================================"

  if [[ ${level} == "HEAD" ]]; then
    line="${colours[DEFAULT]}${msg}${colours[DEFAULT]}"
  else
    line="${colours[${level}]}${date} [${level}] ${msg}${colours[DEFAULT]}"
  fi

  case "${level}" in
  "INFO" | "WARN")
    echo -e "${line}"
    ;;
  "DEBUG")
    if [ "${debug}" -gt 0 ]; then
      echo -e "${line}"
    fi
    ;;
  "ERROR")
    echo -e "${line}" >&2
    ;;
  "HEAD")
    echo -e "${div}\n${line}\n${div}"
    ;;
  esac
}
