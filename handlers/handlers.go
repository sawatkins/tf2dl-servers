package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sawatkins/upfast-tf/database"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).Render("404", fiber.Map{
		"Message": "404 Not found! Please try again",
	}, "layouts/main")
}

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":               "upfast.tf",
		"Canonical":           "https://upfast.tf",
		"Robots":              "noindex, nofollow",
		"Description":         "upfast.tf is a simple tf2 server hosting service",
		"Keywords":            "upfast.tf, upfast, tf2, servers, hosting, game, server, hosting",
		"TotalPlayerSessions": database.GetTotalPlayerSessions(),
		"TotalTimePlayed":     database.GetTotalTimePlayed(),
	}, "layouts/main")
}

func About(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title":       "About - upfast.tf",
		"Canonical":   "https://upfast.tf/about",
		"Robots":      "noindex, nofollow",
		"Description": "About upfast.tf",
		"Keywords":    "upfast.tf, upfast, tf2, servers, hosting, game, server, hosting",
	}, "layouts/main")
}

func GetServerIPs(c *fiber.Ctx) error {
	ips, err := database.GetServerIPs()
	if err != nil {
		return c.Status(500).SendString("Error getting server ips")
	}
	return c.Status(200).JSON(ips)
}

func GetServerInfo(c *fiber.Ctx) error {
	ip := c.Query("ip")
	if ip == "" {
		return c.Status(400).SendString("Missing IP query parameter")
	}

	serverInfo, err := database.GetServerInfo(ip)
	if err != nil {
		c.Status(500).SendString("Error getting server info")
	}

	return c.Status(200).JSON(serverInfo)
}
