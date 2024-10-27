package database

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorcon/rcon"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sawatkins/upfast-tf/models"
)

var db *sql.DB

func InitDB(filepath string) {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connected")
}

func InitServerTable() {
	createServerTableSQL := `
	CREATE TABLE IF NOT EXISTS servers (
		instance_id VARCHAR(20) PRIMARY KEY,
		public_ip CHAR(15),
		public_dns VARCHAR(100),
		name VARCHAR(50),
		server_hostname VARCHAR(100),
		map VARCHAR(50),
		players INTEGER,
		max_players INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	executeSQL(createServerTableSQL)

	log.Println("Servers table created")
}

func InitPlayerSessionTable() {
	createPlayerSessionTableSQL := `
	CREATE TABLE IF NOT EXISTS player_sessions (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	steam_id TEXT NOT NULL,
    	connect_time TEXT NOT NULL,
    	disconnect_time TEXT,
    	duration INTEGER,
		public_ip CHAR(15)
	);`

	executeSQL(createPlayerSessionTableSQL)

	log.Println("PlayerSession table created")
}

func executeSQL(sqlStatement string) {
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Error executing SQL statement: %v", err)
	}
}

func WriteServerToDB(server *models.Server) error {
	writeServerSQL := `
	INSERT INTO servers (
		instance_id, public_ip, public_dns, name, server_hostname, map, players, max_players, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	statement, err := db.Prepare(writeServerSQL)
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		server.InstanceID,
		server.PublicIP,
		server.PublicDNS,
		server.Name,
		server.ServerHostname,
		server.Map,
		server.Players,
		server.MaxPlayers,
		server.CreatedAt,
	)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	log.Println("Server record inserted")
	return nil
}

func GetServerIPs() ([]string, error) {
	rows, err := db.Query("SELECT public_ip FROM servers")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var ips []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		ips = append(ips, ip)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return ips, nil
}

func GetServerInfo(ip string) (models.ServerStatus, error) {
	query := `
	SELECT public_ip, map, players, max_players, server_hostname
	FROM servers
	WHERE public_ip = ?;`

	var serverStatus models.ServerStatus
	err := db.QueryRow(query, ip).Scan(
		&serverStatus.PublicIP,
		&serverStatus.Map,
		&serverStatus.Players,
		&serverStatus.MaxPlayers,
		&serverStatus.Hostname,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No server found with IP: %s", ip)
			return serverStatus, nil
		}
		log.Printf("Error querying server info for IP %s: %v", ip, err)
		return serverStatus, err
	}

	return serverStatus, nil
}

func GetTotalPlayerSessions() int {
	var count int
	query := "SELECT COUNT(*) FROM player_sessions"

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error querying total player sessions: %v", err)
		return 0
	}

	return count
}

func GetTotalTimePlayed() int {
	var totalDuration int
	query := "SELECT SUM(duration) FROM player_sessions"

	err := db.QueryRow(query).Scan(&totalDuration)
	if err != nil {
		log.Printf("Error querying total time played: %v", err)
		return 0
	}

	return totalDuration / 60 // return minues
}
// UpdateServerInfo updates the server information and active player connection in the db for each server IP
func UpdateServerInfo(prevPlayerConnections *map[string]map[string]int64) {
	ips, err := GetServerIPs()
	if err != nil {
		log.Printf("Error updating server info: %v", err)
		return
	}

	for _, ip := range ips {
		rconPass := os.Getenv("RCON_PASSWORD")

		client, err := rcon.Dial(ip+":27015", rconPass)
		if err != nil {
			log.Printf("Failed to connect to RCON: %v", err)
			return
			// TODO: delete prev sessions, or continue loop, if this fails
		}
		defer client.Close()

		response, err := client.Execute("status")
		if err != nil {
			log.Printf("Failed to execute RCON command: %v", err)
			return
		}

		// get server status
		hostname := extract(`hostname:\s*(.+)`, response)
		gameMap := extract(`map\s*:\s*([^\s]+)`, response)
		players, maxPlayers := extractPlayers(`players\s*:\s*(\d+)\s*humans.*\((\d+)\s*max\)`, response)

		// Update the server information in the database
		serverStatus := models.ServerStatus{
			PublicIP:   ip,
			Map:        gameMap,
			Players:    players,
			MaxPlayers: maxPlayers,
			Hostname:   hostname,
		}

		updateServerSQL := `
		UPDATE servers
		SET map = ?, players = ?, max_players = ?, server_hostname = ?
		WHERE public_ip = ?;`

		statement, err := db.Prepare(updateServerSQL)
		if err != nil {
			log.Printf("Error preparing SQL statement: %v", err)
			return
		}
		defer statement.Close()

		_, err = statement.Exec(
			serverStatus.Map,
			serverStatus.Players,
			serverStatus.MaxPlayers,
			serverStatus.Hostname,
			serverStatus.PublicIP,
		)
		if err != nil {
			log.Printf("Error executing SQL statement: %v", err)
			return
		}

		log.Printf("Server info updated for IP: %s", ip)


		// Update active player connections
		currentPlayerIds := extractUniqueIDs(response)
		// log.Println("currentPlayerIds", currentPlayerIds)

		// get new ids (ids in current players not in prev ids)
		for _, currID := range currentPlayerIds {
			if _, exists := (*prevPlayerConnections)[ip][currID]; !exists {
				// for newIds, create new entry in prevplayer connections
				if (*prevPlayerConnections)[ip] == nil {
					(*prevPlayerConnections)[ip] = make(map[string]int64)
				}
				(*prevPlayerConnections)[ip][currID] = time.Now().Unix()
			}
		}

		// log.Println("prevPlayerConenctions", (*prevPlayerConnections))

		// get disconnedted ids (ids in prev ids not in current players)
		disconnectedIds := []string{}
		for prevID := range (*prevPlayerConnections)[ip] {
			if !contains(currentPlayerIds, prevID) {
				disconnectedIds = append(disconnectedIds, prevID)
			}
		}

		// log.Println("disconnectedIds", disconnectedIds)

		// for disconnectedIds, add player session to the db
		for _, id := range disconnectedIds {
			connectTime := time.Unix((*prevPlayerConnections)[ip][id], 0)
			disconnectTime := time.Now()
			duration := disconnectTime.Sub(connectTime)
			newPlayerSession := models.PlayerSession {
				SteamID: id,
				ConnectTime: connectTime.String(),
				DisconnectTime: disconnectTime.String(),
				Duration: int(duration.Seconds()),
				PublicIP: ip,
			}

			// add newPlayerSession to the player_sessions table
			insertPlayerSessionSQL := `
			INSERT INTO player_sessions (
				steam_id, connect_time, disconnect_time, duration, public_ip
			) VALUES (?, ?, ?, ?, ?);`

			statement, err := db.Prepare(insertPlayerSessionSQL)
			if err != nil {
				log.Printf("Error preparing SQL statement for player session: %v", err)
				continue
			}
			defer statement.Close()

			_, err = statement.Exec(
				newPlayerSession.SteamID,
				newPlayerSession.ConnectTime,
				newPlayerSession.DisconnectTime,
				newPlayerSession.Duration,
				newPlayerSession.PublicIP,
			)
			if err != nil {
				log.Printf("Error executing SQL statement for player session: %v", err)
				continue
			}

			log.Printf("Player session recorded for SteamID: %s", id)

			// delete entry from prevPlayerConnections
			delete((*prevPlayerConnections)[ip], id)
		}

	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extract(pattern, response string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(response)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractPlayers(pattern, response string) (string, string) {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(response)
	if len(match) > 2 {
		return match[1], match[2]
	}
	return string(rune(0)), string(rune(0))
}

func extractUniqueIDs(response string) []string {
	re := regexp.MustCompile(`U:1:\d+`)
	matches := re.FindAllString(response, -1)
	return matches
}
