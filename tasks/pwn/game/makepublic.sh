#!/bin/bash

set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

cd deploy/
docker build -f build.Dockerfile -t build-game .
docker run --name build-game-1 build-game
docker cp build-game-1:/build/vuln .
docker rm build-game-1
docker rmi build-game

cd $curdir
mkdir -p $pubtemp/game
cp deploy/vuln $pubtemp/game
cp deploy/Dockerfile $pubtemp/game

cd $pubtemp
zip -r game.zip game

cd $curdir
mv $pubtemp/game.zip public/

rm -rf $pubtemp
