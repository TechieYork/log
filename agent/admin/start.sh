#!/bin/bash

app_name=dm_log_agent

echo "Starting $app_name ... "

echo "Current directory:"$(dirname $(readlink -f $0))

cd `dirname $0`

rm /var/tmp/dark_metrix_log.sock

../bin/dm_log_agent &

sleep 1

chown nobody:nobody /var/tmp/dark_metrix_log.sock

chmod 777 /var/tmp/dark_metrix_log.sock

echo "Already started $app_name"
