#!/bin/bash

cd .

cmake .

make clean; make -j2

mv ./libnet_log_ext.so ./net_log_ext.so
