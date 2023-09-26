#!/bin/bash

set -e

curdir=$(pwd)
pubtemp=$(mktemp -d)

echo $pubtemp

mkdir -p $pubtemp/gif0day
cp -r deploy/* $pubtemp/gif0day/

echo "FLAG=brics+{exampleflag}" > $pubtemp/gif0day/server.env

cd $pubtemp
zip -r "ppc-gif0day.zip" gif0day

cd $curdir

mv $pubtemp/ppc-gif0day.zip public/

rm -rf $pubtemp