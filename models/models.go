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

type ServerStatus struct {
	PublicIP   string `json:"public_ip"`
	Map        string `json:"map"`
	Players    string `json:"players"`
	MaxPlayers string `json:"max_players"`
	Hostname   string `json:"hostname"`
}

type PlayerSession struct {
	SteamID        string `json:"steam_id"`
	ConnectTime    string `json:"connect_time"`
	DisconnectTime string `json:"disconnect_time,omitempty"`
	Duration       int    `json:"duration"` // seconds
	PublicIP       string `json:"public_ip"`
}
