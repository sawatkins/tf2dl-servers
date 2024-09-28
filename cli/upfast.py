import argparse
import subprocess
import os
import sys
import shutil
import tomllib
import json

def get_server_config_name():
    print("Available server configs:")
    for config in os.listdir("./server-configs"):
        if config != "base.toml":
            print(f"  {config}")
    server_config_name = input("Enter name of server config to use: ")
    if server_config_name not in os.listdir("./server-configs"):
        print(f"Error: Server config '{server_config_name}' not found.")
        sys.exit(1)
    return server_config_name

def read_server_config_file(server_config):
    with open(f"./server-configs/{server_config}", "rb") as f:
        config = tomllib.load(f)
    return config

def write_server_to_curent_servers_file(new_server):
    servers_file = "./current-servers.json"
    
    if os.path.exists(servers_file):
        with open(servers_file, "r") as f:
            current_servers = json.load(f)
    else:
        current_servers = {}
    
    current_servers[new_server["instance_id"]] = {
        "public_ip": new_server["public_ip"],
        "public_dns": new_server["public_dns"]
    }
    
    with open(servers_file, "w") as f:
        json.dump(current_servers, f, indent=4)

def create_server():
    server_config_name = get_server_config_name()
    server_config = read_server_config_file(server_config_name)
    base_config = read_server_config_file("base.toml")
    
    # create server with terraform
    try:
        subprocess.run([
            "terraform", "apply",
            "-var", f"region={server_config['region']}",
            "-var", f"instance_type={server_config['instance_type']}",
            "-var", f"key_name={base_config['key_name']}",
            "-var", f"security_group_id={base_config['security_group_id']}",
            "-var", f"iam_instance_profile={base_config['iam_instance_profile']}"
        ], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Error: Terraform apply failed with exit code {e.returncode}")
        sys.exit(1)
    
    # save server info to current-servers.json
    new_server = {
        "instance_id": subprocess.check_output(["terraform", "output", "-raw", "instance_id"]).decode().strip(),
        "public_ip": subprocess.check_output(["terraform", "output", "-raw", "instance_public_ip"]).decode().strip(),
        "public_dns": subprocess.check_output(["terraform", "output", "-raw", "instance_public_dns"]).decode().strip()
    }
    write_server_to_curent_servers_file(new_server)
    
    
    

def check_dependencies():
    required_programs = ["aws", "ansible", "terraform"]
    for program in required_programs:
        if not shutil.which(program):
            print(f"Error: {program} is not installed or not found in PATH.")
            sys.exit(1)

def main():
    check_dependencies()
    os.chdir(os.path.dirname(os.path.abspath(__file__)))

    parser = argparse.ArgumentParser(description="manage upfast.tf servers")
    parser.add_argument("command", choices=["create"], help="command to execute")

    args = parser.parse_args()

    if args.command == "create":
        create_server()
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == "__main__":
    main()