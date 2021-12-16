#!/bin/bash
set -e
BASE_PATH=$(
  cd "$(dirname "$0")" || exit
  pwd
)
CODE_PATH=${BASE_PATH}/..
DOCKER_IMAGE='gate-way'

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

echo "Executing protocol buffer code and swagger generation commands: [5 in total]"

GOOGLE_PB_PKG_PATH="${CODE_PATH}/api"
# alias to docker cmd in docker-build mode
if [ -z "${docker_build}" ]; then
  echo "[Using local tools]"
else
  if [ "${docker_build}" = "yes" ]; then
    if [ -z "$(docker images -q ${DOCKER_IMAGE})" ]; then
      echo "${DOCKER_IMAGE} not exist"
      docker pull ${DOCKER_IMAGE}
    fi
    alias protoc='docker run -it -v ${CODE_PATH}:/data -w /data ${DOCKER_IMAGE} protoc'
    alias find='docker run -it -v ${CODE_PATH}:/data -w /data ${DOCKER_IMAGE} find'
    # use latest googleapis in mod path or in api/googleapis
    GOOGLE_PB_PKG_PATH="/data/api"
    echo "[Using tools in docker]"
  else
    echo "[Using local tools]"
  fi
fi

# shellcheck disable=SC2044
for p in $(find "${GOOGLE_PB_PKG_PATH}/server" -name "*.proto"); do
  echo $p
  PB_GENERATE_CMD="protoc -I. -I api -I ${GOOGLE_PB_PKG_PATH}/googleapis --go_out=. --go-grpc_out=. --proto_path=${GOOGLE_PB_PKG_PATH}/server $p"
  PB_PROXY_GENERATE_CMD="protoc -I. -I api -I ${GOOGLE_PB_PKG_PATH}/googleapis  --grpc-gateway_out=logtostderr=true:. --proto_path=${GOOGLE_PB_PKG_PATH}/server $p"
  echo "$PB_GENERATE_CMD"
  eval "${PB_GENERATE_CMD}"
  echo "$PB_PROXY_GENERATE_CMD"
  eval "${PB_PROXY_GENERATE_CMD}"
done

# shellcheck disable=SC2044
for p in $(find "${GOOGLE_PB_PKG_PATH}/service" -name "*.proto"); do
  echo $p
  PB_GENERATE_CMD="protoc -I. -I api -I ${GOOGLE_PB_PKG_PATH}/googleapis --go_out=. --go-grpc_out=. --proto_path=${GOOGLE_PB_PKG_PATH}/service  $p"
  PB_PROXY_GENERATE_CMD="protoc -I. -I api -I ${GOOGLE_PB_PKG_PATH}/googleapis  --grpc-gateway_out=logtostderr=true:. --proto_path=${GOOGLE_PB_PKG_PATH}/service $p"
  echo "$PB_GENERATE_CMD"
  eval "${PB_GENERATE_CMD}"
  echo "$PB_PROXY_GENERATE_CMD"
  eval "${PB_PROXY_GENERATE_CMD}"
done
