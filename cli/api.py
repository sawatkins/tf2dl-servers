from fastapi import FastAPI
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware  # Add this import
import docker # type: ignore
import re

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allows all origins
    allow_credentials=True,
    allow_methods=["*"],  # Allows all methods
    allow_headers=["*"],  # Allows all headers
)

@app.get("/")
async def root():
    return JSONResponse({"message": "hello root"})

@app.get("/server-info")
async def server_info():
    client = docker.from_env()
    
    try:
        container = client.containers.get('tf2-dedicated')
        
        # Get environment variables from the container
        env_dict = container.attrs['Config']['Env']
        env_dict = dict(item.split('=', 1) for item in env_dict)
        
        server_dir = env_dict.get('SERVER_DIR')
        ip = env_dict.get('IP')
        port = env_dict.get('PORT')
        rcon_password = env_dict.get('RCON_PASSWORD')

        # Run RCON command
        rcon_command = f"{server_dir}/rcon -H {ip} -p {port} -P {rcon_password} status"
        output = container.exec_run(rcon_command).output.decode('utf-8')
        
        # Parse the output
        public_ip = re.search(r'public IP from Steam: (\d+\.\d+\.\d+\.\d+)', output)
        map_name = re.search(r'map\s+:\s+(\w+)', output)
        players = re.search(r'players\s+:\s+(\d+)\s+humans,\s+\d+\s+bots\s+\((\d+)\s+max\)', output)
        hostname = re.search(r'hostname:\s*(.+)', output)

        return JSONResponse({
            "public_ip": public_ip.group(1) if public_ip else None,
            "map": map_name.group(1) if map_name else None,
            "players": int(players.group(1)) if players else None,
            "max_players": int(players.group(2)) if players else None,
            "hostname": hostname.group(1).strip() if hostname else None
        })
    except docker.errors.NotFound:
        return JSONResponse({"error": "Container not found"}, status_code=404)
    except docker.errors.APIError as e:
        return JSONResponse({"error": f"Docker API error: {str(e)}"}, status_code=500)
    except Exception as e:
        return JSONResponse({"error": f"An error occurred: {str(e)}"}, status_code=500)

if __name__ == "__main__":
    import uvicorn
    print("Starting API")
    uvicorn.run(app, host="0.0.0.0", port=8000)
