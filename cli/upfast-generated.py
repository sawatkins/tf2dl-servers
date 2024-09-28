#!/usr/bin/env python3

import argparse
import subprocess
import os
import sys

def get_public_dns():
    if os.path.isfile("hosts.ini"):
        with open("hosts.ini") as f:
            for line in f:
                if "[tf2_server]" in line:
                    return next(f).split()[0]
    else:
        print("Error: hosts.ini not found. Has the server been created?", file=sys.stderr)
        sys.exit(1)

def create_server():
    if subprocess.call(["terraform", "apply"]) != 0:
        print("Error: Terraform apply failed or was declined.")
        sys.exit(1)

    public_ip = subprocess.check_output(["terraform", "output", "-raw", "instance_public_ip"]).decode().strip()
    public_dns = subprocess.check_output(["terraform", "output", "-raw", "instance_public_dns"]).decode().strip()

    hosts_content = f"""[tf2_server]
{public_dns} ansible_user=ec2-user ansible_ssh_private_key_file={os.getenv('SSH_PRIVATE_KEY_PATH')}

[tf2_server:vars]
rcon_password={os.getenv('RCON_PASSWORD')}
server_hostname={os.getenv('SERVER_HOSTNAME')}
"""

    try:
        with open("hosts.ini", "w") as f:
            f.write(hosts_content)
    except IOError:
        print("Error: Failed to create hosts.ini file.")
        sys.exit(1)

    subprocess.call(["sleep", "5"])
    if subprocess.call(["ansible", "-i", "hosts.ini", "-m", "ping", public_dns, "-e", "ansible_ssh_common_args='-o StrictHostKeyChecking=no'"]) != 0:
        print("Error: Ansible ping failed. Check server reachability and configuration.")
        sys.exit(1)

    print("Server is reachable and configured correctly.")
    print(f"IP: {public_ip}")
    print("\nRun 'upfast.py start' to start the server")

def start_server():
    if subprocess.call(["ansible-playbook", "-i", "hosts.ini", "tf2_server_playbook.yml"]) == 0:
        public_ip = subprocess.check_output(["terraform", "output", "-raw", "instance_public_ip"]).decode().strip()
        subprocess.call(["aws", "s3", "cp", "s3://upfast-tf2-hosts/servers.txt", "servers.txt"])
        with open("servers.txt", "a+") as f:
            f.seek(0)
            if public_ip not in f.read():
                f.write(f"{public_ip}\n")
                subprocess.call(["aws", "s3", "cp", "servers.txt", "s3://upfast-tf2-hosts/servers.txt"])
                print("Server started successfully and IP address added to S3 file.")
            else:
                print("Server started successfully. IP address already exists in S3 file.")
        os.remove("servers.txt")
    else:
        print("Error: Failed to start the server.")
        sys.exit(1)

def restart_server():
    subprocess.call(["ansible-playbook", "-i", "hosts.ini", "tf2_server_playbook.yml", "--tags", "restart"])

def destroy_server():
    public_ip = subprocess.check_output(["terraform", "output", "-raw", "instance_public_ip"]).decode().strip()

    if subprocess.call(["terraform", "destroy"]) == 0:
        subprocess.call(["aws", "s3", "cp", "s3://upfast-tf2-hosts/servers.txt", "servers.txt"])
        with open("servers.txt", "r") as f:
            lines = f.readlines()
        with open("servers.txt", "w") as f:
            for line in lines:
                if public_ip not in line:
                    f.write(line)
        subprocess.call(["aws", "s3", "cp", "servers.txt", "s3://upfast-tf2-hosts/servers.txt"])
        os.remove("servers.txt")
        os.remove("hosts.ini")
        print("Server destroyed, IP removed from S3 file, and hosts.ini cleared.")
    else:
        print("Error: Failed to destroy server.")
        sys.exit(1)

def connect_to_server():
    public_dns = get_public_dns()
    subprocess.call(["ssh", "-i", os.getenv("SSH_PRIVATE_KEY_PATH"), f"ec2-user@{public_dns}"])

def connect_to_container():
    public_dns = get_public_dns()
    subprocess.call(["ssh", "-i", os.getenv("SSH_PRIVATE_KEY_PATH"), f"ec2-user@{public_dns}", "-t", "sudo docker attach tf2-dedicated"])

def main():
    parser = argparse.ArgumentParser(description="Manage TF2 server")
    parser.add_argument("command", choices=["create", "start", "restart", "destroy", "connect", "connect_container"], help="Command to execute")

    args = parser.parse_args()

    if args.command == "create":
        create_server()
    elif args.command == "start":
        start_server()
    elif args.command == "restart":
        restart_server()
    elif args.command == "destroy":
        destroy_server()
    elif args.command == "connect":
        connect_to_server()
    elif args.command == "connect_container":
        connect_to_container()
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == "__main__":
    main()