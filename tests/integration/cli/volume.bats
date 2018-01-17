#!/usr/bin/env ./bats-0.4.0/bin/bats
# load bats-0.4.0/test/test_helper
#
# bats test for the volume

@test "sloppy start app with invalid volume size" {
  run sloppy start command/testdata/invalid_testproject_with_volume.json
  echo $output
  [[ "$output" =~  "error: Invalid volume size for /usr/share/docs Volume size needs to be a multiple of 8GB, was 7GB" ]]
}

@test "sloppy start app with a volume size > as the bought addon" {
  run sloppy start command/testdata/testproject_with_too_big_volume.json
  echo $output
  [[ "${lines[0]}" =~  "error: Could not validate input Overrun maximum storage." ]]
}

@test "sloppy start app with valid volume" {
  run sloppy start command/testdata/testproject_with_volume.json
  echo $output
}

@test "sloppy show to see if volume exists" {
  run sloppy show volumetest/frontend/apache
  result=$(echo ${lines[7]} | tr -d [:space:])
  echo $result
  [[ $result ==  "Volumes:'/usr/share/docs'8GB" ]]
}

@test "sloppy delete volumetest" {
  run sloppy delete volumetest
  echo $output
  [ "$status" -eq 0 ]
  [[ "${lines[0]}" =~ "Project volumetest successfully deleted" ]]
}
