#!/usr/bin/env bash


POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -h|--help)
    HELP=yes
    shift # past argument
    ;;
    --update-snapshots)
    UPDATE_SNAPSHOTS=yes
    shift # past argument
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

function help {
  cat << EOF
Test runner
Runs tests for the Bluelink GitHub Registry:
bash scripts/run-tests.sh
EOF
}

if [ -n "$HELP" ]; then
  help
  exit 0
fi

set -e
echo "" > coverage.txt

# Ensure test env vars are exported for tests to use.
set -a
. .env.test
set +a

go test -timeout 30000ms -race -coverprofile=coverage.txt -coverpkg=./... -covermode=atomic `go list ./... | egrep -v '(/(testutils|tools))$'`

if [ -z "$GITHUB_ACTION" ]; then
  # We are on a dev machine so produce html output of coverage
  # to get a visual to better reveal uncovered lines.
  go tool cover -html=coverage.txt -o coverage.html
fi

if [ -n "$GITHUB_ACTION" ]; then
  # We are in a CI environment so run tests again to generate JSON report.
  go test -timeout 30000ms -json -tags "$TEST_TYPES" `go list ./... | egrep -v '(/(testutils|tools))$'` > report.json
fi
