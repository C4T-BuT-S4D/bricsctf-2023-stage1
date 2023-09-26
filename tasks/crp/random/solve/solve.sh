#!/bin/bash
g++ gen_rels.cpp -Ofast -march=native -lm4ri -o a.out
g++ hax.cpp -Ofast -march=native -lm4ri -o b.out
./a.out > list_rels
xxd -r -p < ../public/output.txt > out_bin
./b.out
