#!/bin/bash
VERSION=${1}
PUSH=${2}
if [ "${VERSION}" != "" ]; then
  docker build -t shreddedbacon/fake-powerwall:arm32v6-rpi-${VERSION} .
  if [ "$PUSH" == "push" ]; then
    docker push shreddedbacon/fake-powerwall:arm32v6-rpi-${VERSION}
  fi
fi

