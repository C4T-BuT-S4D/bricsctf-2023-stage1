#!/bin/bash
perl -MO=C,-oscript.c script.pl
perl /usr/local/bin/cc_harness -O2 -o perlovka script.c

