#!/bin/bash
# script to update tf2 and restart the game server services
set -e

echo "switching to gameserver user..."
sudo -u gameserver /home/gameserver/hlserver/tf2_ds

echo "restarting service..."
sudo systemctl restart tf2server.service

echo "success!"
