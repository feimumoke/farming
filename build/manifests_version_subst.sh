#!/bin/bash

set -e
BASE_PATH=$(cd "$(dirname "$0")" || exit;pwd)
MANIFESTS_PATH=${BASE_PATH}/../deployments/manifests

bash "${BASE_PATH}"/envsubst_install.sh

cd "$MANIFESTS_PATH" || exit
for directory in *; do
    if [ -d "$directory" ]; then
        cd "${directory}" || exit
        # shellcheck disable=SC2016
        # this is intended to be single quote so that envsubst knows which environment variable to substitute
        for deploymentFile in *deployment.yaml; do
            tmp="temp_${deploymentFile}"
            APP_VERSION="${VERSION}" envsubst '$APP_VERSION' < "${deploymentFile}" > "${tmp}"
            mv "${tmp}" "${deploymentFile}"
        done
        cd ../
    fi
done