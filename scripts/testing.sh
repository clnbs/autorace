#!/bin/bash
green() {
  "$@" | GREP_COLORS='mt=01;32' grep --color .
}

red() {
  "$@" | GREP_COLORS='mt=01;31' grep --color .
}

yellow() {
  "$@" | GREP_COLORS='mt=01;93' grep --color .
}

check_command_success() {
  CODE_TO_COMPARE_TO=$2
  RETURNED_CODE=$1
  if [ $RETURNED_CODE -ne $CODE_TO_COMPARE_TO ]; then
    if [[ $2 != "" ]]; then
      red echo "$3"
    fi
    exit 1
  fi
}

set_env() {
  export POD_IMAGE_NAME="dummyImage"
  export POD_VERSION="dummy-1.0.0"
  export POD_HOSTNAME="dummy"
  export POD_ENV="ENVIRONMENT=${ENVIRONMENT};LOG_LEVEL=debug"
  export POD_ARGS="ls;-lah"
  export POD_NETWORKS="dummy_net;logs"
}

unset_env() {
  unset POD_IMAGE_NAME
  unset POD_VERSION
  unset POD_HOSTNAME
  unset POD_ENV
  unset POD_ARGS
  unset POD_NETWORKS
}

start_test() {
  set_env
  make run
  red echo "waiting services to be fully started ..."
  sleep 20
  go test -v -cover ./...
  make down
  unset_env
}

exported_environment=false

if [ -z "$ENVIRONMENT" ]; then
  exported_environment=true
  export ENVIRONMENT=dev
fi

yellow echo "starting test for ${ENVIRONMENT} environment"
start_test
green echo "test completed"

if [ "$exported_environment" = true ]; then
  unset ENVIRONMENT
fi


