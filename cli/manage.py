#!/usr/bin/env python3

import argparse
import subprocess
import os
import sys
import shutil
import json
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
    headers = { "Authorization": os.getenv("CLI_AUTH_KEY").strip('"') }
    port = 8080 #(5000) TODO: find better way to manage this
    
    for instance_id, server_info in current_servers.items():
        payload = {
            "instance_id": instance_id,
            "public_ip": server_info.get("public_ip"),
            "public_dns": server_info.get("public_dns"),
            "name": server_info.get("name"),
            "server_hostname": server_info.get("server_hostname"),
        }
        
        try:
            response = requests.post(f"http://localhost:{port}/api/current-servers", json=payload, headers=headers)
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

# create_server creates a all servers with terraform
def create_server():
    try:
        subprocess.run([
            "terraform", "apply",
            "-var-file", f"./terraform.tfvars"
        ], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Error: Terraform apply failed with exit code {e.returncode}")
        sys.exit(1)
    
    tf2_server_us = {
        "instance_id": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_us_id"]).decode().strip(),
        "public_ip": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_us_public_ip"]).decode().strip(), # TODO get the elastic ip
        "public_dns": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_us_public_dns"]).decode().strip(),
        "name": "tf2_server_us",
        "server_hostname": "simple surf server (us) - servers.tf2dl.net"
    }
    write_server_to_curent_servers_file(tf2_server_us)
    
    tf2_server_eu = {
        "instance_id": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_eu_id"]).decode().strip(),
        "public_ip": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_eu_public_ip"]).decode().strip(),
        "public_dns": subprocess.check_output(["terraform", "output", "-raw", "tf2_server_eu_public_dns"]).decode().strip(),
        "name": "tf2_server_eu",
        "server_hostname": "simple surf server (eu) - servers.tf2dl.net"
    }
    write_server_to_curent_servers_file(tf2_server_eu)
    
    post_current_servers_to_db()
    
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

def connect_to_server():
    print_current_servers()
    server_name = input("Enter the server name to connect to (Name): ").strip()
    
    current_servers = read_current_servers_file()
    for server_info in current_servers.values():
        if server_info["name"] == server_name:
            public_ip = server_info["public_ip"]
            try:
                subprocess.run(["ssh", "-i", os.getenv("SSH_PRIVATE_KEY_PATH").strip('"'), f"admin@{public_ip}"], check=True)
            except subprocess.CalledProcessError as e:
                print(f"Error: SSH connection failed with exit code {e.returncode}")
            return
    
    print(f"No server found with the name: {server_name}")
        
def destroy_server():
    # for now, delete all
    subprocess.run(["terraform", "destroy", "-var-file", f"./upfast.tfvars"], check=True)
    os.remove("./current-servers.json")

def check_dependencies():
    required_programs = ["aws", "terraform"]
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
            if not line.startswith("#"):
                key, value = line.strip().split("=")
                os.environ[key] = value

    parser = argparse.ArgumentParser(description="manage servers.tf2dl.net servers")
    parser.add_argument("command", choices=["create", "destroy", "list", "connect", "write_db"], help="command to execute")

    args = parser.parse_args()
    if args.command == "create":
        create_server()
    elif args.command == "destroy":
        destroy_server()
    elif args.command == "list":
        print_current_servers()
    elif args.command == "connect":
        connect_to_server()
    elif args.command == "write_db":
        post_current_servers_to_db()
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == "__main__":
    main()