#!/bin/bash
set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

cp -R deploy $pubtemp/picoblog
cd $pubtemp

sed -i -r 's/brics\+\{.+\}/brics+{fake_flag}/g' picoblog/blog/docker-compose.yml
zip -9 -r picoblog.zip picoblog

cd $curdir
mv $pubtemp/picoblog.zip public
rm -rf $pubtemp
