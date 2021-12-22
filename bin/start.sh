#!/bin/bash
BASE_PATH=$(
    cd "$(dirname "$0")" || exit
    pwd
)
if [ -z "${DOMAIN}" ]; then
    # shellcheck disable=SC2046
    # take advantage of xargs construct arguments separated by space
    export $(grep -E -v '^#' "${BASE_PATH}"/../.env | xargs)
else
    # shellcheck disable=SC2046
    # take advantage of xargs construct arguments separated by space
    export $(grep -E -v '^#' "${BASE_PATH}"/../.env | grep -v DOMAIN | xargs)
fi
SERVICE=${PROJECT_NAME}
BIN_PATH=${BASE_PATH}/../bin
ROOT_PATH=${BASE_PATH}/..

chmod +x "${BIN_PATH}"/*
cd "${ROOT_PATH}" || exit

sync

export PROJECT_NAME=$PROJECT_NAME
export GRPC_GO_LOG_SEVERITY_LEVEL=info
export GRPC_GO_LOG_VERBOSITY_LEVEL=99

GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn ./bin/"${SERVICE}"
