#!/bin/sh
BASE_PATH=$(cd "$(dirname "$0")" || exit;pwd)
if [ -z "${DOMAIN}" ]
then
    # shellcheck disable=SC2046
    # take advantage of xargs construct arguments separated by space
    export $(grep -E -v '^#' "${BASE_PATH}"/../.env | xargs)
else
    # shellcheck disable=SC2046
    # take advantage of xargs construct arguments separated by space
    export $(grep -E -v '^#' "${BASE_PATH}"/../.env | grep -v DOMAIN | xargs)
fi
SERVICE=${PROJECT_NAME}

pid=$(ps -ef | grep "${SERVICE}" | grep -v "grep" | grep -v "bin/start" | awk '{print $2}')
if [ -z "${pid}" ]
then
    echo "${SERVICE} not exist"
else
    kill -INT "${pid}"
 	echo "send signal to stop ${SERVICE}"
fi