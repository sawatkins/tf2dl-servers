#!/usr/bin/env python3

import argparse
import subprocess
import os
import sys
import shutil
import json
import time
# import tempfile
import requests


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

def post_current_servers_to_db():
    current_servers = read_current_servers_file()
    headers = { "Authorization": os.getenv("CLI_AUTH_KEY") }
    
    for instance_id, server_info in current_servers.items():
        payload = {
            "instance_id": instance_id,
            "public_ip": server_info.get("public_ip"),
            "public_dns": server_info.get("public_dns"),
            "name": server_info.get("name"),
            "server_hostname": server_info.get("server_hostname"),
        }
        
        try:
            response = requests.post("http://localhost:80/api/current-servers", json=payload, headers=headers)
            response.raise_for_status()  
            print(f"Posted server {instance_id} to database.")
        except requests.exceptions.RequestException as e:
            print(f"Error posting server {instance_id} to database: {e}")

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
    
# def push_current_servers_to_s3():
#     current_servers = read_current_servers_file()
    
#     content = '\n'.join(server['public_ip'] for server in current_servers.values())
    
#     with tempfile.NamedTemporaryFile(mode='w+', delete=False) as temp_file:
#         temp_file.write(content)
#         temp_file_path = temp_file.name
    
#     try:
#         subprocess.run([
#             "aws", "s3", "cp",
#             temp_file_path,
#             "s3://upfast-tf2-hosts/servers.txt"
#         ], check=True)
#         print("Successfully pushed server IPs to S3")
#     except subprocess.CalledProcessError as e:
#         print(f"Error pushing server IPs to S3: {e}")
#     finally:
#         os.unlink(temp_file_path)

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
    
    post_current_servers_to_db()
    # push_current_servers_to_s3()
    
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
    # for now, delete all
    subprocess.run(["terraform", "destroy", "-var-file", f"./upfast.tfvars"], check=True)
    os.remove("./inventory.ini")
    os.remove("./current-servers.json")

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