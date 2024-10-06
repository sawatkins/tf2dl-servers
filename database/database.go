package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	log.Println("Table created")
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

// UpdateServerInfo updates the server information in the database for each server IP
func UpdateServerInfo() {
	ips, err := GetServerIPs()
	if err != nil {
		log.Printf("Error updating server info: %v", err)
		return
	}

	for _, ip := range ips {
		resp, err := http.Get(fmt.Sprintf("http://%s:8000/server-info", ip))
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Error getting server info for %s, status %d: %v", ip, resp.StatusCode, err)
			return
		}
		defer resp.Body.Close()

		var serverStatus models.ServerStatus
		if err := json.NewDecoder(resp.Body).Decode(&serverStatus); err != nil {
			log.Printf("Error parsing JSON response from server: %v", err)
			return
		}

		// Update the server information in the database
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
	}
}
