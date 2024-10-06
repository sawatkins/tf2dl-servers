package main

import (
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"github.com/sawatkins/upfast-tf/database"
	"github.com/sawatkins/upfast-tf/handlers"
)

func main() {
	port := flag.String("port", ":8080", "Port to listen on")
	dev := flag.Bool("dev", true, "Enable development mode")
	flag.Parse()

	if err := godotenv.Load("./cli/.env"); err != nil {
		log.Fatalln("Did not load .env file")
	}

	database.InitDB("./data/upfast.db")
	database.InitServerTable()
	go startServerInfoUpdater()

	engine := html.New("./templates", ".html")
	if *dev {
		engine.Reload(true)
		engine.Debug(true)
	}

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Static("/", "./static")

	app.Post("/api/current-servers", handlers.PostCurrentServer)
	app.Get("/api/server-ips", handlers.GetServerIPs)
	app.Get("/api/server-info", handlers.GetServerInfo)

	app.Get("/", handlers.Index)
	app.Get("/about", handlers.About)
	app.Use(handlers.NotFound)

	log.Println("Server starting on port", *port)
	log.Fatal(app.Listen(*port)) // default port: 8080
}

func startServerInfoUpdater() {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		database.UpdateServerInfo()
	}
}
