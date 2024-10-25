#!/bin/bash

# This script installs a Team Fortress 2 server following this guide
# https://wiki.teamfortress.com/wiki/Linux_dedicated_server
# for Debian/Ubuntu Linux machines

set -e

check_prerequisites() {
    echo "Checking prerequisites..."
    if ! grep -qEi "(debian|ubuntu)" /etc/os-release; then
        echo "This script is intended for Debian/Ubuntu systems only."
        exit 1
    fi

    local required_space=15000000 # 15GB in KB
    local available_space=$(df / | tail -1 | awk '{print $4}')
    if [ "$available_space" -lt "$required_space" ]; then
        echo "Not enough disk space. At least 15GB is required."
        exit 1
    fi
}

# CHECK PREREQUISITES
check_prerequisites

# INSTALL REQUIREMENTS
echo "Setting up system requirements..."
sudo dpkg --add-architecture i386
sudo apt-get update
sudo apt-get install -y lib32z1 libncurses5:i386 libbz2-1.0:i386 lib32gcc-s1 lib32stdc++6 libtinfo5:i386 libcurl3-gnutls:i386 libsdl2-2.0-0:i386

# SETUP USER & INSTALL STEAMCMD
echo "Creating gameserver user and setting up SteamCMD..."
sudo adduser --disabled-login --gecos "" gameserver
sudo -u gameserver bash -c "
    mkdir -p /home/gameserver/hlserver &&
    cd /home/gameserver/hlserver &&
    wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz &&
    tar zxf steamcmd_linux.tar.gz &&
    rm steamcmd_linux.tar.gz
"

# DOWNLOAD THE SERVER
echo "Downloading Team Fortress 2 server..."
sudo -u gameserver bash -c "
    /home/gameserver/hlserver/steamcmd.sh +force_install_dir /home/gameserver/hlserver/tf2 +login anonymous +app_update 232250 +quit
"

# CREATE & UPDATE SERVER CONFIG
echo "Setting up server configuration files..."
sudo -u gameserver bash -c "
    cd /home/gameserver/hlserver/tf2/tf/cfg &&
    touch server.cfg motd.txt mapcycle.txt &&
    echo 'Please configure the files server.cfg, motd.txt, mapcycle.txt'
"

# TODO: download custom servers from s3

echo "Installation complete."