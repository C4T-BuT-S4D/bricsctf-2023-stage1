#!/bin/bash

echo "--> Creating AVD"
./create-avd.sh

echo "--> Starting AVD from the snapshot"
emulator @default -cores 8 -no-window -sdcard ./sdcard.img -snapshot droidchat &

echo "--> Waiting 10 seconds for emulator boot"
sleep 10

echo "--> Starting runner"
./runner droidchat

printf --  "--> Runner exited with code %d. Finishing emulator" $?
kill -INT %1
wait %1

