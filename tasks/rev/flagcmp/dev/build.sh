#!/bin/sh
set -e

docker build -t flagcmp-wasm:latest .

docker run -it -v ./pkg:/pkg flagcmp-wasm cp -r /build/pkg /
