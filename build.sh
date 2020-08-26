#!/bin/bash

COMMIT_ID=$(git rev-parse HEAD)
APP_VERSION=$(git tag --contains "${COMMIT_ID}")
APP_COMMIT=$(git rev-parse --short HEAD)
APP_CREATED=$(date '+%Y-%m-%dT%H:%M:%SZ%Z')
if [ -z "${APP_VERSION}" ]
then
      APP_VERSION="${COMMIT_ID}"
fi

go build \
	    -mod=vendor \
	    -ldflags "\
	        -w \
          -X github.com/spyzhov/healthy/app.Version=${APP_VERSION} \
	        -X github.com/spyzhov/healthy/app.Commit=${APP_COMMIT} \
	        -X github.com/spyzhov/healthy/app.Created=${APP_CREATED} \
	    " \
	    -o "$2" "$1"
