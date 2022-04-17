#!/usr/bin/env bash

REPO_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
package_name=`cat $REPO_DIR/NAME`

(set -x; go build $package_name.go && mv $package_name $REPO_DIR/bin/$package_name)
