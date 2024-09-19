package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).Render("404", fiber.Map{
		"Message": "404 Not found! Please try again",
	}, "layouts/main")
}

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":       "upfast.tf",
		"Canonical":   "https://upfast.tf",
		"Robots":      "noindex, nofollow",
		"Description": "upfast.tf is a simple tf2 server hosting service",
		"Keywords":    "upfast.tf, upfast, tf2, servers, hosting, game, server, hosting",
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
	const awsEndpoint = "https://bwdfgz2pbedm7ficoxqxbhfazi0ynfoh.lambda-url.us-west-1.on.aws"

	resp, err := http.Get(awsEndpoint)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Failed to fetch server IPs: %v", err))
	}
	defer resp.Body.Close()

	return c.Status(resp.StatusCode).JSON(resp.Body)
}

func GetServerInfo(c *fiber.Ctx) error {
	ip := c.Query("ip")
	if ip == "" {
		return c.Status(400).SendString("Missing IP query parameter")
	}

	url := fmt.Sprintf("http://%s:8000/server-info", ip)
	resp, err := http.Get(url)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Failed to fetch server info: %v", err))
	}
	defer resp.Body.Close()

	return c.Status(resp.StatusCode).JSON(resp.Body)
}