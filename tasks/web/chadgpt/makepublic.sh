#!/bin/bash

set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

echo $pubtemp

mkdir -p $pubtemp/chadgpt
cp -r deploy/* $pubtemp/chadgpt/

sed -i '' -r 's/brics\+{.*}/brics+{fake}/g' $pubtemp/chadgpt/conf/db_init/db.sql

cd $pubtemp

zip -r "web-chadgpt.zip" chadgpt

cd $curdir

mv $pubtemp/web-chadgpt.zip public/

rm -rf $pubtemp
