#!/bin/bash

set -e
BASE_PATH=$(cd "$(dirname "$0")" || exit;pwd)
MANIFESTS_BASE_PATH=${BASE_PATH}/../deployments/base

if [ "$#" -ne 2 ] ; then
  echo "Usage: $0 docker.io/test/test TAG"
  exit 1
fi

bash "${BASE_PATH}/kustomize_install.sh"

cd "${MANIFESTS_BASE_PATH}" || exit
kustomize edit set image "$1"="$1":"$2"

