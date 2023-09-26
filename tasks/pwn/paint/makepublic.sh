#!/bin/bash

set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

cd deploy/
docker build -f build.Dockerfile -t build-paint .
docker run --name build-paint-1 build-paint
docker cp build-paint-1:/build/vuln .
docker rm build-paint-1
docker rmi build-paint

cd $curdir
mkdir -p $pubtemp/paint
cp deploy/vuln $pubtemp/paint
cp deploy/Dockerfile $pubtemp/paint

cd $pubtemp
zip -r paint.zip paint

cd $curdir
mv $pubtemp/paint.zip public/

rm -rf $pubtemp
