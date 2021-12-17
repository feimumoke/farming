#!/bin/bash

set -e
current_directory=$PWD

if ! [ -x "$(command -v envsubst)" ]; then
    echo "Prepare to install envsubst."
    if [ -f "${current_directory}"/envsubst ]; then
        echo "A file named envsubst already exists in current working directory. Please remove it first."
        exit 1
    fi
    echo "Downloading envsubst to current directory."
    if [[ "$OSTYPE" == linux* || "$OSTYPE" == darwin* ]]; then
        echo "Exec curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-$(uname -s)-$(uname -m) -o envsubst"
        curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-$(uname -s)-$(uname -m) -o envsubst

        if [ -e "envsubst" ]; then
            chmod +x envsubst
            sudo mv envsubst /usr/local/bin || exit
            if [ -x "$(command -v kustomize)" ]; then
                echo "envsubst successfully installed."
            else
                echo "Failed to install envsubst. Please move envsubst to bin manually."
            fi
        else
            echo "Failed to download envsubst."
            exit 1
        fi
    else
        echo "Not supported OS. Failed to install envsubst."
        exit 1
    fi
else
    echo "envsubst has already been installed."
fi
