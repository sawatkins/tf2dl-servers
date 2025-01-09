#!/bin/bash

# Script to install a dedicated Team Fortress 2 Linux server
# - Checks prerequisites
# - Installs required packages
# - Creates gameserver user
# - Installs TF2 server files
# - Downloads custom maps

set -e

# Make sure script is run as root
if [ "$EUID" -ne 0 ]; then
  echo "This script must be run as root."
  exit 1
fi

# Check prerequisites
if ! grep -qx 'ID=debian' /etc/os-release; then
    echo "This script is for Debian only."
    exit 1
fi

required_space=15000000 # 15GB
available_space=$(df / | tail -1 | awk '{print $4}')
if [ "$available_space" -lt "$required_space" ]; then
    echo "Not enough disk space. At least 15GB is required."
    exit 1
fi

# Install required packages
dpkg --add-architecture i386
apt-get update
apt-get install -y \
  lib32z1 \
  libncurses5:i386 \
  libbz2-1.0:i386 \
  lib32gcc-s1 \
  lib32stdc++6 \
  libtinfo5:i386 \
  libcurl3-gnutls:i386 \
  libsdl2-2.0-0:i386 \
  wget \
  tar

# Create gameserver user
useradd -m -s /bin/bash -U gameserver
passwd -l gameserver

# Install tf2 server as gameserver user
su - gameserver <<'EOF'
set -e

mkdir -p ~/hlserver
cd ~/hlserver

# download and install tf2 server files
wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz
tar -xzf steamcmd_linux.tar.gz
rm steamcmd_linux.tar.gz
./steamcmd.sh +login anonymous +force_install_dir ~/hlserver/tf2 +app_update 232250 validate +quit

# install custom maps
mkdir -p ~/hlserver/tf2/tf/maps
cd ~/hlserver/tf2/tf/maps
wget https://upfast.tf/maps/surf_utopia_v3.bsp
wget https://upfast.tf/maps/surf_pulse_v2.bsp

# create tf2 config files
cd ~/hlserver/tf2/tf/cfg
touch server.cfg motd.txt mapcycle.txt
echo "Please configure server.cfg, motd.txt, and mapcycle.txt as needed."
EOF

echo "Installation complete"
