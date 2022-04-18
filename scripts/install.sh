#!/usr/bin/env bash

export package_name=gofuzzyclone
export author=amazingandyyy
# get the latest version or change to a specific version
export latest_version=$(curl --silent "https://api.github.com/repos/$author/$package_name/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
[[ -n "$1" ]] && latest_version=$1

echo "Installing/upgrading to $package_name@$latest_version (latest)"

fmpfolder=/tmp
tmpoupoutput=$fmpfolder/$package_name-$latest_version
tmpoupoutputgz=$tmpoupoutput.tar.gz
userlocalbin=/usr/local/bin/$package_name

curl -Ls "https://github.com/$author/$package_name/archive/refs/tags/$latest_version.tar.gz" -o $tmpoupoutputgz
tar -zxf $tmpoupoutputgz --directory /tmp
if ! ls -d $userlocalbin > /dev/null 2>&1
then
  sudo touch $userlocalbin
  sudo mv $fmpfolder/$package_name-$latest_version/bin/$package_name $userlocalbin
else
  sudo mv $fmpfolder/$package_name-$latest_version/bin/$package_name $userlocalbin
fi

rm -rf $fmpfolder/$package_name $tmpoupoutput $tmpoupoutputgz

if ! [[ -x $(command -v $package_name) ]]; then
  echo "Installed $package_name unsuccessfully" >&2
  exit 1
else
  echo "Installed $package_name@$latest_version successfully!"
  gofuzzyclone -help
fi
