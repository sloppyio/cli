#!/usr/bin/env ./bats-0.4.0/bin/bats
# load bats-0.4.0/test/test_helper

url_before_change="sloppy-cli-testing.sloppy.zone"
url_after_change="sloppy-cli-after-testing.sloppy.zone"

@test "running just sloppy must return usage" {
  run sloppy
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "usage" ]]
}

@test "running sloppy start without correct parameter" {
  run sloppy start
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "error: 'start' requires a minimum of 1 argument." ]]
}

@test "running sloppy start, file not found" {
  run sloppy start foobarbaz
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "error: file 'foobarbaz' not found." ]]
}

@test "sloppy start testproject with forbidden env SLOPPY_*" {
  env_test_helper "SLOPPY_TEST"
}

@test "sloppy start testproject with forbidden env MARATHON_*" {
  env_test_helper "MARATHON_TEST"
}

@test "sloppy start testproject with forbidden env MESOS_*" {
  env_test_helper "MESOS_TEST"
}

@test "sloppy start testproject with forbidden env WEAVE_CIDR_*" {
  env_test_helper "WEAVE_CIDR_TEST"
}

@test "sloppy logs without correct parameter" {
  run sloppy logs
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "error: 'logs' requires a minimum of 1 argument." ]]
}

@test "sloppy scale without correct parameter" {
  run sloppy scale
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: 'scale' requires a minimum of 2 arguments" ]]
}

@test "sloppy scale without correct instance parameter" {
  run sloppy scale apache/frontend/apache a
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: invalid instance number 'a'" ]]
}

@test "sloppy scale with project and not app" {
  run sloppy scale apache 2
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: invalid application path 'apache'" ]]
}

@test "sloppy restart without correct parameter" {
  run sloppy restart
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: 'restart' requires a minimum of 1 argument." ]]
}

@test "sloppy restart with project and not app" {
  run sloppy restart apache
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: invalid application path 'apache'" ]]
}

@test "sloppy change without correct parameter" {
  run sloppy change
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: 'change' requires a minimum of 1 argument." ]]
}

@test "sloppy change with project and no json file" {
  run sloppy change letschat
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: application path or project file required." ]]
}

@test "sloppy change with no changes" {
  run sloppy change apache/frontend/apache
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: missing options." ]]
}

@test "sloppy start testproject" {
  run sloppy start command/testdata/testproject.json
  echo $output
  [ "$status" -eq 0 ]
  [[ "${lines[3]}" =~  "frontend" ]]
}

@test "sloppy start different project with same domain" {
  run sloppy start command/testdata/testproject_duplicate_domain.json
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: Could not validate input *.domain \"${url_before_change}\" is already used" ]]
}

@test "sloppy show with wrong name" {
  run sloppy show apache1
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: Couldn't find project apache1" ]]
}

@test "sloppy show" {
  sleep 15
  run sloppy show
  echo $output
  [ "$status" -eq 0 ]
  [[ "${lines[3]}" =~  "apache" ]]
}

@test "sloppy get url before change" {
  result=$(http --headers GET $url_before_change | head -n 1 | tr -d '[:space:]')
  echo $result
  [ ${result} == "HTTP/1.1200OK" ]
}

@test "sloppy show apache" {
  run sloppy show apache
  echo $output
  [ "$status" -eq 0 ]
  [[ "${lines[3]}" =~  "frontend" ]]
}

@test "sloppy show apache/frontend" {
  run sloppy show apache/frontend
  echo $output
  [ "$status" -eq 0 ]
  [[ "${lines[3]}" =~  "apache" ]]
}

@test "sloppy show apache/frontend/apache" {
  run sloppy show apache/frontend/apache
  echo $output
  [[ "$output" =~  "apache" ]]
}

@test "sloppy stats" {
  run sloppy stats apache
  echo $output
  echo $status
  [ "$status" -eq 0 ]
  [[ "$output" =~  "CONTAINER" ]]
}

@test "sloppy logs with wrong project name" {
  run sloppy logs xyz
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: Couldn't find project xyz" ]]
}

@test "sloppy logs apache/frontend/apache" {
  run sleep 15 && kill -SIGTERM $(pgrep sloppy) &
  run sloppy logs apache/frontend/apache
  echo $output
  [[ "${lines[0]}" =~  "apache" ]]
}

@test "sloppy change apache/frontend/apache" {
  run sloppy change -e TESTVAR:XYZ apache/frontend/apache
  echo $output
  [ "$status" -eq 0 ]
  [[ "$output" =~ "TESTVAR=\"XYZ\"" ]]
}

@test "sloppy change domain for apache/frontend/apache" {
  run sleep 20
  run sloppy change -d $url_after_change apache/frontend/apache
  echo $output
  [ "$status" -eq 0 ]
  [[ "$output" =~ $url_after_change ]]
}

@test "sloppy get url after change" {
  result=$(http --headers GET $url_after_change | head -n 1 | tr -d '[:space:]')
  echo $result
  [ ${result} == "HTTP/1.1200OK" ]
}

@test "sloppy change env to test rollback for apache/frontend/apache" {
  run sleep 15
  run sloppy change -e test:rollback apache/frontend/apache
  echo $output
  [ "$status" -eq 0 ]
  [[ "$output" =~ $url_after_change ]]
}

@test "sloppy rollback to previous version of apache/frontend/apache" {
  run sleep 15
  version=$(sloppy show apache/frontend/apache | grep -A 1 Versions | tail -n 1 | tr -d [:space:])
  run sloppy rollback apache/frontend/apache $version
  echo $output
  [ "$status" -eq 0 ]
}

@test "Check if rollback works as expected." {
  run sleep 15
  run sloppy show apache/frontend/apache
  echo $output | grep -v "rollback"
}

@test "sloppy scale apache/frontend/apache" {
  skip
  run sleep 15
  run sloppy scale apache/frontend/apache 2
  echo $output
  [ "$status" -eq 0 ]
  [[ $output =~ "Instances: 1 / 2" ]]
}

@test "sloppy delete apache" {
  run sloppy delete apache
  echo $output
  [ "$status" -eq 0 ]
  [[ "${lines[0]}" =~ "Project apache successfully deleted" ]]
}

env_test_helper() {
  local -r env="$1"
  local -r filename=$(mktemp /tmp/XXXXXXXXXXXX.json)
  cat command/testdata/testproject.json | jq ".services[0].apps[0].env.${env}=\"forbidden\"" > $filename
  run sloppy start $filename
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~ "error: Invalid project schema instance.services[0].apps[0].env {\"${env}\":\"forbidden\"} additionalProperty \"${env}\" exists in instance when not allowed" ]]
  rm $filename
}
