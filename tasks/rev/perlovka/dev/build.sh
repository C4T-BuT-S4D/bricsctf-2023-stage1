#!/bin/bash

docker build -t perlovka . 
docker run -it -v $(pwd):/ctf perlovka /bin/bash
