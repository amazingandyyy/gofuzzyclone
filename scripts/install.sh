#!/usr/bin/env bash
export stable_version=0.0.6

export REPO_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
export package_name=$(basename $REPO_DIR)
export author=amazingandyyy
[[ -n "$1" ]] && stable_version=$1

echo "> installing $package_name@$stable_version"
curl -LsO https://github.com/$author/$package_name/archive/refs/tags/$stable_version.zip &&
unzip -o $stable_version.zip &&
rm -rf /opt/homebrew/bin/$package_name &&
sudo touch /opt/homebrew/bin/$package_name &&
chmod +x $package_name-$stable_version/bin/$package_name &&
mv -f $package_name-$stable_version/bin/$package_name /opt/homebrew/bin
rm -rf $package_name-$stable_version $stable_version.zip

if ! [[ -x $(command -v $package_name) ]]; then
  echo 'Error: gofuzzyclone failed to install' >&2
  exit 1
else
  echo "> install $package_name@$stable_version successfully!"
  gofuzzyclone -help
fi
