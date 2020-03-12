#!/bin/bash
cd ../../stardust-experimental/
#make
cd build/
sudo xl create -q ./mini-os.conf
cp mini-os.gz ../../debugger/
cd ../../debugger/build/
gunzip mini-os.gz
sudo xl list
