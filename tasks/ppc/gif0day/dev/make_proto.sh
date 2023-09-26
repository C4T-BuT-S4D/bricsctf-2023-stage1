#!/bin/bash

outDir="$1"
python -m grpc_tools.protoc -I . --python_out="$outDir"  --pyi_out="$outDir" --grpclib_python_out="$outDir" restore.proto