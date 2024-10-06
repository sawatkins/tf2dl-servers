package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
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