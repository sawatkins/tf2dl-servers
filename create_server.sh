#!/bin/bash

source .env

if ! terraform apply; then
    echo "Error: Terraform apply failed or was declined."
    exit 1
fi

# Extract output values
public_ip=$(terraform output -raw instance_public_ip)
public_dns=$(terraform output -raw instance_public_dns)

# Create hosts.ini file
if ! cat > hosts.ini <<-EOF
[tf2_server]
$public_dns ansible_user=ec2-user ansible_ssh_private_key_file=$SSH_PRIVATE_KEY_PATH
EOF
then
    echo "Error: Failed to create hosts.ini file."
    exit 1
fi

if ! ansible -i hosts.ini -m ping $public_dns; then
    echo "Error: Ansible ping failed. Check server reachability and configuration."
    exit 1
fi

echo "Server is reachable and configured correctly."
echo "IP: $public_ip"