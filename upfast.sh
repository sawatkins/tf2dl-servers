#!/bin/bash

source .env

create_server() {
    if ! terraform apply; then
        echo "Error: Terraform apply failed or was declined."
        exit 1
    fi

    # Extract output values
    public_ip=$(terraform output -raw instance_public_ip)
    public_dns=$(terraform output -raw instance_public_dns)

    if ! printf '%s\n' \
        "[tf2_server]" \
        "$public_dns ansible_user=ec2-user ansible_ssh_private_key_file=$SSH_PRIVATE_KEY_PATH" \
        > hosts.ini
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
}

start_server() {
    ansible-playbook -i hosts.ini tf2_.yml
}

create_server