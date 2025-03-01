#!/bin/bash
#
# Script to connect to each remote tf2 server with rcon

if [[ $# -ne 1 || !( "$1" == "tf2_server_us" || "$1" == "tf2_server_eu") ]]; then
    echo "usage: rcon_connect.sh <tf2_server_us|tf2_server_eu>"
    exit 1
fi

cd ~/rcon || exit 1

if [ ! -f "rcon.yaml" ]; then
    echo "error: rcon.yaml file doesn't exist"
    exit 1
fi

./rcon -e "$1"
