from fastapi import FastAPI
from fastapi.responses import JSONResponse
import subprocess
import os
import re

app = FastAPI()

@app.get("/")
async def root():
    return JSONResponse({"message": "hello root"})

@app.get("/server-info")
async def server_info():
    server_dir = os.environ.get('SERVER_DIR', '')
    ip = os.environ.get('IP', '127.0.0.1')
    port = os.environ.get('PORT', '')
    rcon_password = os.environ.get('RCON_PASSWORD', '')

    docker_command = f"sudo docker exec tf2-dedicated {server_dir}/rcon -H {ip} -p {port} -P {rcon_password} status"
    
    try:
        output = subprocess.check_output(docker_command, shell=True, text=True)
        
        # Parse the output
        public_ip = re.search(r'public IP from Steam: (\d+\.\d+\.\d+\.\d+)', output)
        map_name = re.search(r'map\s+:\s+(\w+)', output)
        players = re.search(r'players\s+:\s+(\d+)\s+humans,\s+\d+\s+bots\s+\((\d+)\s+max\)', output)

        return JSONResponse({
            "public_ip": public_ip.group(1) if public_ip else None,
            "map": map_name.group(1) if map_name else None,
            "human_players": int(players.group(1)) if players else None,
            "max_players": int(players.group(2)) if players else None
        })
    except subprocess.CalledProcessError as e:
        return JSONResponse({"error": f"Command execution failed: {str(e)}"}, status_code=500)
    except Exception as e:
        return JSONResponse({"error": f"An error occurred: {str(e)}"}, status_code=500)

if __name__ == "__main__":
    import uvicorn
    print("Starting API")
    uvicorn.run(app, host="0.0.0.0", port=8000)
