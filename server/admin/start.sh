#!/bin/bash

app_name=dm_log_server

echo "Starting $app_name ... "

echo "Current directory:"$(dirname $(readlink -f $0))

cd `dirname $0`

../bin/dm_log_server &

echo "Already started $app_name"
