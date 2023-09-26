#!/bin/bash

set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

echo $pubtemp

mkdir -p $pubtemp/gigachadgpt
cp -r deploy/* $pubtemp/gigachadgpt/

sed -i '' -r 's/brics\+{.*}/brics+{fake}/g' $pubtemp/gigachadgpt/conf/db_init/db.sql

cd $pubtemp

zip -r "web-gigachadgpt.zip" gigachadgpt

cd $curdir

mv $pubtemp/web-gigachadgpt.zip public/

rm -rf $pubtemp
