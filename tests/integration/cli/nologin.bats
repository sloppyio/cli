#!/usr/bin/env ./bats-0.4.0/bin/bats
# load bats-0.4.0/test/test_helper

@test "running sloppy show without login must return Sloppy Error" {
  run sloppy show
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: not logged in" ]]
}

@test "running sloppy show without login must return Sloppy Error" {
  run sloppy show
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: not logged in" ]]
}

@test "sloppy scale not logged in" {
  run sloppy scale apache/frontend/apache 2
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "not logged in" ]]
}

@test "sloppy restart not logged in" {
  run sloppy restart apache/frontend/apache
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: not logged in" ]]
}

@test "sloppy change not logged in" {
  run sloppy change -i 2 apache/frontend/apache
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: not logged in" ]]
}

@test "sloppy stats not logged in" {
  run sloppy stats apache
  echo $output
  [ "$status" -eq 1 ]
  [[ "${lines[0]}" =~  "error: not logged in" ]]
}
