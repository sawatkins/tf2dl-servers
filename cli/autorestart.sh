#!/bin/bash
# Restarts the tf2 server if memory usage is >93% and no players are connected

# needs to run as root
if [ "$EUID" -ne 0 ]; then
    echo "script must be run as superuser"
    exit 1
fi

# make sure jq is installed
if ! command -v jq > /dev/null 2>&1; then
    echo "jq is not installed"
    exit 1
fi

# date for logs
date 

# if service is down, restart
if ! systemctl is-active --quiet tf2server.service; then
  systemctl restart tf2server.service
  echo "service was down, restarted..."
  exit 0
fi

# check memory usage
memory_used=$(free -m | awk '/Mem:/ {print int($3/$2 * 100)}')
if [ "$memory_used" -lt 93 ]; then
    echo "memory usage is less than 93%. exiting..."
    exit 0
fi

# check player count
if [ ! -e "/home/admin/public_ip" ]; then
    curl -s ifconfig.me/ip > /home/admin/public_ip
    echo "created file /home/admin/public_ip"
fi

public_ip=$(cat /home/admin/public_ip)
server_info=$(curl -s "https://upfast.tf/api/server-info?ip=$public_ip")
if [ -z "$server_info" ]; then
    echo "failed to get server info"
    exit 1
fi

players=$(echo "$server_info" | jq -r '.players')
if [ "$players" -gt 0 ]; then
    echo "server is not empty. exiting..."
    exit 0
fi

# restart the server
echo "restarting tf2 service..."
systemctl restart tf2server.service
