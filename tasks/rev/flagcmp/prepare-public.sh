#!/bin/bash
set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

cp -R deploy $pubtemp/flagcmp
cd $pubtemp

sed -i -r 's/brics\+\{.+\}/brics+{fake_flag}/g' flagcmp/docker-compose.yml
zip -9 -r flagcmp.zip flagcmp

cd $curdir
mv $pubtemp/flagcmp.zip public
rm -rf $pubtemp
