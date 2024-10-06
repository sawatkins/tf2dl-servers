package database

import (
	"database/sql"
	"log"

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

