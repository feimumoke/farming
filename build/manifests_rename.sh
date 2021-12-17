#!/bin/bash

set -e
BASE_PATH=$(cd "$(dirname "$0")" || exit;pwd)
MANIFESTS_PATH=${BASE_PATH}/../deployments/manifests

cd "$MANIFESTS_PATH" || exit
# Note: This scripts rename the .yaml files generated by Kustomize under deployments/manifests
# Following part only works perfectly if .yaml files under deployments/base have distinct values for "kind".
# If any two .yaml files contain the same value for "kind", this script will add a count number to the filename.
for directory in *; do
    if [ -d "$directory" ]; then
        # Will not run if no directories are available
        cd "$directory" || exit
        for file in *; do
            # !!!Important!!!
            # Note: the general naming rule of Kustomize is:
            # namespace_api_version_type_application[_]name.yaml
            # e.g. namespace_apps_v1_deployment_application_name.yaml
            # e.g. namespace_monitoring.coreos.com_v1_servicemonitor_appname.yaml
            NEW_FILE_PATH=$(basename "$file" '.yaml' | awk -F"_" '{if ($1 ~/v[0-9]+*/) print$3"_"$2;else print$4"_"$3}')
            if [ -f "$NEW_FILE_PATH.yaml" ]; then
                COUNT=1
                while [ -f "$NEW_FILE_PATH$COUNT.yaml" ]; do
                    COUNT=$((COUNT+1))
                done
                mv "$file" "$NEW_FILE_PATH$COUNT.yaml"
            else
                mv "$file" "$NEW_FILE_PATH.yaml"
            fi
        done
        cd ..
    fi
done