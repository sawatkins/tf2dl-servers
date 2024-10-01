#!/usr/bin/env python3

import argparse
import subprocess
import os
import sys
import shutil
import json
import time


def write_server_to_curent_servers_file(new_server):
    current_servers = read_current_servers_file()
    current_servers[new_server["instance_id"]] = {
        "public_ip": new_server["public_ip"],
        "public_dns": new_server["public_dns"],
        "name": new_server["name"],
        "server_hostname": new_server["server_hostname"]
    }

    with open("./current-servers.json", "w") as f:
        json.dump(current_servers, f, indent=4)

def read_current_servers_file():
    if os.path.exists("./current-servers.json"):
        with open("./current-servers.json", "r") as f:
            return json.load(f)
    else:
        return {}
    
def update_ansible_inventory():
    with open("./inventory.ini", "w") as f:
        f.write("[tf2_server]\n")
        for _, server_info in read_current_servers_file().items():
            f.write(f"{server_info['public_dns']} server_hostname='{server_info['server_hostname']}'\n")
        f.write("[tf2_server:vars]\n")
        f.write("ansible_user=ec2-user\n")
        f.write(f"ansible_ssh_private_key_file={os.getenv('SSH_PRIVATE_KEY_PATH')}\n")
        f.write(f"rcon_password={os.getenv('RCON_PASSWORD')}\n")

# create_server creates a all servers with terraform
def create_server():
    try:
        subprocess.run([
            "terraform", "apply",
            "-var-file", f"./upfast.tfvars"
        ], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Error: Terraform apply failed with exit code {e.returncode}")
        sys.exit(1)
    
    # save server info to current-servers.json
    tf2_server_jump_01 = {
        "instance_id": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_jump_01_id"]).decode().strip(),
        "public_ip": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_jump_01_public_ip"]).decode().strip(),
        "public_dns": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_jump_01_public_dns"]).decode().strip(),
        "name": "tf2_server_jump_01",
        "server_hostname": "jump 24/7 - upfast.tf"
    }
    write_server_to_curent_servers_file(tf2_server_jump_01)
    
    tf2_server_surf_01 = {
        "instance_id": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_surf_01_id"]).decode().strip(),
        "public_ip": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_surf_01_public_ip"]).decode().strip(),
        "public_dns": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_surf_01_public_dns"]).decode().strip(),
        "name": "tf2_server_surf_01",
        "server_hostname": "surf 24/7 - upfast.tf"
    }
    write_server_to_curent_servers_file(tf2_server_surf_01)
    
    print("updating ansible inventory")
    update_ansible_inventory()
    
    time.sleep(5)
    
    print("running ansible playbook")
    try:
        subprocess.run([
            "ansible-playbook",
            "tf2_server_playbook.yml",
            "-i", "./inventory.ini",
            "-e", "ansible_ssh_common_args='-o StrictHostKeyChecking=no'"
        ], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Error: Ansible playbook failed with exit code {e.returncode}")
        sys.exit(1)
    
def print_current_servers():
    current_servers = read_current_servers_file()
    if not current_servers:
        print("No current servers.")
        return
    print("Current servers:")
    for instance_id, server_info in current_servers.items():
        print(f"Instance ID: {instance_id}")
        print(f"  Public IP: {server_info['public_ip']}")
        print(f"  Public DNS: {server_info['public_dns']}")
        print(f"  Name: {server_info['name']}")
        print(f"  Server Hostname: {server_info['server_hostname']}")
        print("")
        
def destroy_server():
    pass
    # current_servers = read_current_servers_file()
    # if not current_servers:
    #     print("No servers to destroy.")
    #     return
    
    # print_current_servers()
    # server_id_to_destroy = input("Enter instance id of server to destroy: ")
    # if server_id_to_destroy not in current_servers:
    #     print(f"Error: Server '{server_id_to_destroy}' not found.")
    #     sys.exit(1)

    # # destroy server with terraform
    # try:
    #     subprocess.run([
    #         "terraform", "destroy",
    #         "-var", f"instance_id={server_id_to_destroy}"
    #     ], check=True)
    # except subprocess.CalledProcessError as e:
    #     print(f"Error: Terraform destroy failed with exit code {e.returncode}")
    #     sys.exit(1)
    
    # # remove server info from current-servers.json
    # del current_servers[server_id_to_destroy]
    # # write_current_servers_(current_servers)

def check_dependencies():
    required_programs = ["aws", "ansible", "terraform"]
    for program in required_programs:
        if not shutil.which(program):
            print(f"Error: {program} is not installed or not found in PATH.")
            sys.exit(1)

def main():
    check_dependencies()
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    
    # read .env file
    with open(".env", "r") as f:
        for line in f:
            key, value = line.strip().split("=")
            os.environ[key] = value

    parser = argparse.ArgumentParser(description="manage upfast.tf servers")
    parser.add_argument("command", choices=["create", "destroy", "list"], help="command to execute")

    args = parser.parse_args()
    if args.command == "create":
        create_server()
    elif args.command == "destroy":
        destroy_server()
    elif args.command == "list":
        print_current_servers()
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == "__main__":
    main()