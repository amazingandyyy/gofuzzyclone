#!/usr/bin/env bash

export version=$(cat VERSION)
export author=amazingandyyy
export package_name=$(cat NAME)

if [ -n "$1" ]; then
  version=$1
fi

echo "> installing package_name@$version"
curl -LsO https://github.com/$author/package_name/archive/refs/tags/$version.zip &&
unzip -o $version.zip &&
rm -rf /opt/homebrew/bin/package_name &&
sudo touch /opt/homebrew/bin/package_name &&
chmod +x package_name-$version/bin/package_name &&
mv -f package_name-$version/bin/package_name /opt/homebrew/bin
rm -rf package_name-$version $version.zip

if ! [ -x "$(command -v package_name)" ]; then
  echo 'Error: package_name failed to install' >&2
  exit 1
else
  echo "> install package_name@$version successfully!"
  eval($package_name -help)
fi
