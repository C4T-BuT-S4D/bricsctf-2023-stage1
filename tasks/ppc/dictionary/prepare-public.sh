#!/bin/bash
set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

cp -R deploy $pubtemp/dictionary
cd $pubtemp

sed -i -r 's/brics\+\{.+\}/brics+{fake_flag}/g' dictionary/docker-compose.yml
zip -9 -r dictionary.zip dictionary

cd $curdir
mv $pubtemp/dictionary.zip public
rm -rf $pubtemp