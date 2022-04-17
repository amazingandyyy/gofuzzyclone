#!/usr/bin/env bash

export package_name=$(cat NAME)

REPO_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"

set -x

go build $package_name.go && mv $package_name $REPO_DIR/bin/$package_name

set +x
