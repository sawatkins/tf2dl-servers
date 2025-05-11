package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"

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
	database.InitPlayerSessionTable()
	go startServerInfoUpdater()
	go checkForGameUpdate()

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
	prevPlayerConnections := map[string]map[string]int64{} // map[ip]map[playerID]timestamp{}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		database.UpdateServerInfo(&prevPlayerConnections)
	}
}

func checkForGameUpdate() {
	prevItemDate := time.Time{}
	ticker := time.NewTicker(3 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		resp, err := http.Get("https://www.teamfortress.com/rss.xml")
		log.Println("Fetching rss feed")
		if err != nil {
			log.Printf("Error fetching TF2 RSS feed: %v", err)
			continue
		}

		feed, err := gofeed.NewParser().Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error parsing TF2 RSS feed: %v", err)
			continue
		}

		if len(feed.Items) == 0 {
			continue
		}

		latestItem := feed.Items[0]
		newItemDate := *latestItem.PublishedParsed
		if !newItemDate.After(prevItemDate) {
			continue
		}

		prevItemDate = newItemDate
		if strings.Contains(latestItem.Title, "Team Fortress 2 Update Released") {
			log.Printf("New TF2 update")
			http.Post(os.Getenv("NOTIFY_URL"), "text/plain", strings.NewReader("New TF2 Update Released!"))
		}
	}
}
