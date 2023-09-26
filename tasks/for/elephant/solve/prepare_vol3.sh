#!/bin/bash

git_clone() {
    git clone --depth=1 https://github.com/$1.git
}

echo "--> Cloning repositories"

git_clone volatilityfoundation/volatility3
git_clone volatilityfoundation/dwarf2json

echo "--> Building dwarf2json"
(
    cd dwarf2json
    go build .
)

echo "--> Downloading vmlinux w/debug symbols (will take up about 6 GiB of disk space)"
mkdir deb
cd deb
wget "http://ddebs.ubuntu.com/ubuntu/pool/main/l/linux/linux-image-unsigned-5.15.0-84-generic-dbgsym_5.15.0-84.93_amd64.ddeb" -O pkg.ddeb
ar x pkg.ddeb
unxz data.tar.xz -c | tar xvf - ./usr/lib/debug/boot/
mv ./usr/lib/debug/boot/* ../vmlinux
cd ..
rm -rf deb

echo "--> Generating volatility symbols"
mkdir -p volatility3/volatility3/symbols/linux
dwarf2json/dwarf2json linux --elf vmlinux > volatility3/volatility3/symbols/linux/5.15.0-84-generic.json

echo "--> Done. Keeping vmlinux"

