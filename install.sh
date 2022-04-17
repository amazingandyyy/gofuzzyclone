#!/usr/bin/env bash

export version=0.0.4
export author=amazingandyyy
export package_name=gofuzzyclone

if [ -n "$1" ]; then
  version=$1
fi

echo "> installing $package_name@$version"
curl -LsO https://github.com/$author/$package_name/archive/refs/tags/$version.zip &&
unzip -o $version.zip &&
rm -rf /opt/homebrew/bin/$package_name &&
sudo touch /opt/homebrew/bin/$package_name &&
chmod +x $package_name-$version/bin/$package_name &&
mv -f $package_name-$version/bin/$package_name /opt/homebrew/bin
rm -rf $package_name-$version $version.zip

if ! [ -x "$(command -v gofuzzyclone)" ]; then
  echo 'Error: gofuzzyclone failed to install' >&2
  exit 1
else
  echo "> install $package_name@$version successfully!"
  gofuzzyclone -help
fi
