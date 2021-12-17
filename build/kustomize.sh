#!/bin/bash

set -e
BASE_PATH=$(
    cd "$(dirname "$0")" || exit
    pwd
)
CODE_PATH=${BASE_PATH}/..
OVERLAYS_PATH=${BASE_PATH}/../deployments/overlays
MANIFESTS_PATH=${BASE_PATH}/../deployments/manifests
DOCKER_IMAGE='gw-kt'
OVERLAYS_PATH_DOCKER="deployments/overlays"
MANIFESTS_PATH_DOCKER="deployments/manifests"

while getopts "d:" arg; do
    # shellcheck disable=SC2220
    case ${arg} in
    d)
        docker_build=${OPTARG}
        ;;
    esac
done

echo "Executing kustomize manifests generation"

if [ -z "${docker_build}" ]; then
    echo "[Using local kustomize]"
    bash "${BASE_PATH}"/kustomize_install.sh
else
    if [ "${docker_build}" = "yes" ]; then
        if [ -z "$(docker images -q "${DOCKER_IMAGE}")" ]; then
            echo "${DOCKER_IMAGE} not exist"
            docker pull "${DOCKER_IMAGE}"
        fi
        alias kustomize='docker run --rm -it -v ${CODE_PATH}:/data -w /data ${DOCKER_IMAGE} kustomize'
        # use latest googleapis in mod path or in api/googleapis
        echo "[Using tools in docker]"
     else
        echo "[Using local tools]"
    fi

fi

cd "$OVERLAYS_PATH" || exit 1
for directory in *; do
    if [ -d "$directory" ]; then
        # Will not run if no directories are available
        if [ -d "$MANIFESTS_PATH"/"$directory" ]; then
            rm -rf "${MANIFESTS_PATH:?}"/"$directory"
        fi
        mkdir -p "$MANIFESTS_PATH"/"$directory"
        if [ "${docker_build}" = "yes" ]; then
            kustomize build "${OVERLAYS_PATH_DOCKER}/${directory}" -o "${MANIFESTS_PATH_DOCKER}/${directory}" || exit 1
        else
            kustomize build "${OVERLAYS_PATH}"/"${directory}" -o "${MANIFESTS_PATH}"/"${directory}" || exit 1
        fi
    fi
done