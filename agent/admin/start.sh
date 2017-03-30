#!/bin/bash

app_name=dm_log_agent

echo "Starting $app_name ... "

echo "Current directory:"$(dirname $(readlink -f $0))

cd `dirname $0`

../bin/dm_log_agent &

echo "Already started $app_name"
