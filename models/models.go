package models

type Server struct {
    InstanceID     string `json:"instance_id"`
    PublicIP       string `json:"public_ip"`
    PublicDNS      string `json:"public_dns"`
    Name           string `json:"name"`
    ServerHostname string `json:"server_hostname"`
    Map            string `json:"map"`
    Players        int    `json:"players"`
    MaxPlayers     int    `json:"max_players"`
    CreatedAt      string `json:"created_at"` 
}