#!/bin/bash

echo "--> Creating AVD default"
avdmanager create avd -n "default" -d 47 -k "system-images;android-33;google_apis_playstore;x86_64" -f
rm sdcard.img*
mksdcard 256M sdcard.img

echo "--> Cold booting AVD"
emulator @default -cores 8 -no-window -sdcard ./sdcard.img &

echo "--> Waiting 30 seconds to boot"
sleep 30

echo "--> Installing app.apk"
adb install app.apk  # mount as a volume!

echo "--> Launching app once to get token"
adb shell monkey -p ru.bricsctf.droidchat -c android.intent.category.LAUNCHER 1

sleep 2
adb shell am force-stop ru.bricsctf.droidchat
sleep 2

echo "--> Taking a snapshot..."
adb emu avd snapshot save "droidchat"

echo "--> Stopping emulator"
kill -INT %1
wait %1
echo "--> Snapshot ready!"


