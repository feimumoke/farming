#!/bin/bash

set -e
current_directory=$PWD

if ! [ -x "$(command -v kustomize)" ]; then
    echo "Prepare to install Kustomize."
    if [ -f "${current_directory}"/kustomize ]; then
        echo "A file named kustomize already exists in current working directory. Please remove it first."
        exit 1
    fi
    os=windows
    if [[ "$OSTYPE" == linux* ]]; then
        os=linux
    elif [[ "$OSTYPE" == darwin* ]]; then
        os=darwin
    else
        echo "Not supported OS. Failed to install Kustomize"
        exit 1
    fi
    echo "Downloading Kustomize to current directory "
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
    if [ -f "kustomize" ]; then
        chmod +x kustomize
        bin=/usr/local/bin/
        mv kustomize "${bin}" || exit
        if [ -x "$(command -v kustomize)" ]; then
            echo "Kustomize successfully installed."
        else
            echo "Failed to install kustomize to ${bin}. Please move kustomize to bin manually."
        fi
    else
        echo "Failed to download kustomize_${os}_amd64.tar.gz from COS."
        exit 1
    fi
else
    echo "Kustomize has already been installed."
fi
