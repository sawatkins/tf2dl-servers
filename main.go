package main

import (
	"flag"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	"github.com/sawatkins/upfast-tf2-web/handlers"
)

func main() {
	port := flag.String("port", ":8080", "Port to listen on")
	prefork := flag.Bool("prefork", false, "Enable prefork in Production")
	dev := flag.Bool("dev", true, "Enable development mode")
	flag.Parse()


	engine := html.New("./templates", ".html")
	if *dev {
		engine.Reload(true)
		engine.Debug(true)
	}

	app := fiber.New(fiber.Config{
		Prefork: *prefork,
		Views:   engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Static("/", "./static")

	app.Get("/api/server-ips", handlers.GetServerIPs)
	app.Get("/api/server-info", handlers.GetServerInfo)

	app.Get("/", handlers.Index)
	app.Get("/about", handlers.About)
	app.Use(handlers.NotFound)

	log.Println("Server starting on port", *port)
	log.Fatal(app.Listen(*port)) // default port: 8080
}
