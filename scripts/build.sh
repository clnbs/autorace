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

build_linux() {
  green echo "Starting building Autorace for linux"
  yellow pwd
  docker build -t autorace_linux:build . -f ./build/package/client/Dockerfile.build.linux
  RESULT=$?
  check_command_success $RESULT 0 "Could not build Autorace for linux"
  docker container create --name extract autorace_linux:build
  RESULT=$?
  check_command_success $RESULT 0 "Could not start builder container"
  docker container cp extract:/go/src/github.com/clnbs/autorace/autorace.bin ./autorace.bin
  RESULT=$?
  check_command_success $RESULT 0 "Could not extract binary from builder image"
  docker container rm -f extract
  RESULT=$?
  check_command_success $RESULT 0 "Could not remove builder container"
}

build_linux_test() {
  green echo "Starting building Autorace for linux"
  yellow pwd
  docker build -t autorace_linux_test:build . -f ./build/package/client/Dockerfile.build.linux_test
  RESULT=$?
  check_command_success $RESULT 0 "Could not build Autorace for linux"
  docker container create --name extract autorace_linux_test:build
  RESULT=$?
  check_command_success $RESULT 0 "Could not start builder container"
  docker container cp extract:/go/src/github.com/clnbs/autorace/autorace_test.bin ./autorace_test.bin
  RESULT=$?
  check_command_success $RESULT 0 "Could not extract binary from builder image"
  docker container rm -f extract
  RESULT=$?
  check_command_success $RESULT 0 "Could not remove builder container"
}

build_windows() {
  green echo "Starting building Autorace for Windows"
  docker build -t autorace_windows:build . -f ./build/package/client/Dockerfile.build.windows
  RESULT=$?
  check_command_success $RESULT 0 "Could not build Autorace for linux"
  docker container create --name extract autorace_windows:build
  RESULT=$?
  check_command_success $RESULT 0 "Could not start builder container"
  docker container cp extract:/go/src/github.com/clnbs/autorace/autorace.exe ./autorace.exe
  RESULT=$?
  check_command_success $RESULT 0 "Could not extract binary from builder image"
  docker container rm -f extract
  RESULT=$?
  check_command_success $RESULT 0 "Could not remove builder container"
}

build_server_image() {
  green echo "Starting building servers images"
  docker build -t autorace_static . -f ./build/package/static/Dockerfile
  RESULT=$?
  check_command_success $RESULT 0 "Could not build static Autorace server"
  docker build -t autorace_dynamic . -f ./build/package/dynamic/Dockerfile
  RESULT=$?
  check_command_success $RESULT 0 "Could not build dynamic Autorace server"
}


clean_linux() {
  docker rmi autorace_linux:build > /dev/null 2>&1
}

clean_linux_test() {
  docker rmi autorace_linux_test:build > /dev/null 2>&1
}
clean_windows() {
  docker rmi autorace_windows:build > /dev/null 2>&1
}

clean_running_dynamic_container() {
  docker rm -f $(docker ps -aq -f name=dynamic) > /dev/null 2>&1
}

down_server_stack() {
  docker-compose --file deployments/docker-compose.yaml down
}

help() {
  echo -e "\nUsage : $0 OPTION"
  echo -e "\nOption :"
  echo -e "\linux\tCompile Autorace client binary for Linux systems"
  echo -e "\windows\tCompile Autorace client binary for Windows systems"
  echo -e "\server\tCompile Autorace server stack"
}

OPTION=$1
DISTRIB=$2

if [ -z "$OPTION"  ]; then
  yellow echo "building all"
  build_server_image
  build_linux
  build_linux_test
  build_windows
  exit 0
elif [[ "$OPTION" == "linux" ]]; then
  yellow echo "building linux"
  build_linux
  exit 0
elif [[ "$OPTION" == "linux_test" ]]; then
  yellow echo "building linux test"
  build_linux_test
  exit 0
elif [[ "$OPTION" == "windows" ]]; then
  yellow echo "building windows"
  build_windows
  exit 0
elif [[ "$OPTION" == "server" ]]; then
  yellow echo "building server"
  build_server_image
  exit 0
elif [[ "$OPTION" == "clean" ]]; then
  if [ -z "$DISTRIB"  ]; then
    yellow echo "cleaning all images"
    clean_linux
    clean_linux_test
    clean_windows
    docker rmi $(docker images -q --filter "dangling=true") > /dev/null 2>&1
  elif [[ "$DISTRIB" == "linux" ]]; then
    yellow echo "cleaning linux image"
    clean_linux
    docker rmi $(docker images -q --filter "dangling=true") > /dev/null 2>&1
    elif [[ "$DISTRIB" == "linux_test" ]]; then
    yellow echo "cleaning linux image"
    clean_linux_test
    docker rmi $(docker images -q --filter "dangling=true") > /dev/null 2>&1
  elif [[ "$DISTRIB" == "windows" ]]; then
    yellow echo "cleaning windows image"
    clean_windows
    docker rmi $(docker images -q --filter "dangling=true") > /dev/null 2>&1
  elif [[ "$DISTRIB" == "server" ]]; then
    yellow echo "cleaning server images"
    docker rmi $(docker images -q --filter "dangling=true") > /dev/null 2>&1
  fi
elif [[ "$OPTION" == "down" ]]; then
  clean_running_dynamic_container
  down_server_stack
fi

