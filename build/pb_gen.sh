#!/bin/bash
set -e
BASE_PATH=$(
  cd "$(dirname "$0")" || exit
  pwd
)
CODE_PATH=${BASE_PATH}/..
DOCKER_IMAGE='feimumoke/go-gate-way:v1'
#DOCKER_IMAGE='my-gt'

# shellcheck disable=SC2046
# take advantage of xargs construct arguments separated by space
export $(grep -E -v '^#' "${BASE_PATH}"/../.env | xargs)

while getopts "d:" arg; do
  # shellcheck disable=SC2220
  case ${arg} in
  d)
    docker_build=${OPTARG}
    ;;
  esac
done

echo "Executing protocol buffer commands"

GOOGLE_PB_PKG_PATH="${CODE_PATH}/api"
SERVER_PB_PATH="${GOOGLE_PB_PKG_PATH}/server/*.proto"
SERVICE_PB_PATH="${GOOGLE_PB_PKG_PATH}/service/*.proto"
# alias to docker cmd in docker-build mode
if [ -z "${docker_build}" ]; then
  echo "[Using local tools]"
else
  if [ "${docker_build}" = "yes" ]; then
    if [ -z "$(docker images -q ${DOCKER_IMAGE})" ]; then
      echo "${DOCKER_IMAGE} not exist"
      docker pull ${DOCKER_IMAGE}
    fi
    alias protoc='docker run --rm -it -v ${CODE_PATH}:/data -w /data ${DOCKER_IMAGE} sh -c "protoc'
    alias find='docker run -it -v ${CODE_PATH}:/data -w /data ${DOCKER_IMAGE} find'
    # use latest googleapis in mod path or in api/googleapis
    GOOGLE_PB_PKG_PATH="/data/api"
    SERVER_PB_PATH="${GOOGLE_PB_PKG_PATH}/server/*.proto\""
    SERVICE_PB_PATH="${GOOGLE_PB_PKG_PATH}/service/*.proto\""
    echo "[Using tools in docker]"
  else
    echo "[Using local tools]"
  fi
fi

PB_GENERATE_CMD="protoc -I ${GOOGLE_PB_PKG_PATH}/googleapis --go_out=. --go-grpc_out=. --proto_path=${GOOGLE_PB_PKG_PATH}/server  ${SERVER_PB_PATH}"
eval "${PB_GENERATE_CMD}"
PB_PROXY_GENERATE_CMD="protoc -I ${GOOGLE_PB_PKG_PATH}/googleapis  --grpc-gateway_out=logtostderr=true:. --proto_path=${GOOGLE_PB_PKG_PATH}/server  ${SERVER_PB_PATH}"
eval "${PB_PROXY_GENERATE_CMD}"
PB_GENERATE_CMD="protoc -I ${GOOGLE_PB_PKG_PATH}/googleapis --go_out=. --go-grpc_out=. --proto_path=${GOOGLE_PB_PKG_PATH}/service  ${SERVICE_PB_PATH}"
eval "${PB_GENERATE_CMD}"
PB_PROXY_GENERATE_CMD="protoc -I ${GOOGLE_PB_PKG_PATH}/googleapis  --grpc-gateway_out=logtostderr=true:. --proto_path=${GOOGLE_PB_PKG_PATH}/service   ${SERVICE_PB_PATH}"
eval "${PB_PROXY_GENERATE_CMD}"

echo "SUCCESS"
