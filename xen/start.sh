#!/bin/bash
cd ../../stardust-experimental/build/
sudo xl create -q ./mini-os.conf
cp mini-os.gz ../../debugger/
cd ../../debugger/
gunzip mini-os.gz
sudo xl list
