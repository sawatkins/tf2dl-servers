package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/sawatkins/tf2dl-servers/database"
	"github.com/sawatkins/tf2dl-servers/models"
)

func PostCurrentServer(c *fiber.Ctx) error {
	expectedKey := os.Getenv("CLI_AUTH_KEY")
	requestKey := c.Get("Authorization")

	if expectedKey != requestKey {
		return c.Status(401).SendString("Unauthorized")
	}

	var newServer models.Server
	if err := c.BodyParser(&newServer); err != nil {
		return c.Status(400).SendString("Bad Request: " + err.Error())
	}

	if err := database.WriteServerToDB(&newServer); err != nil {
		return c.Status(500).SendString("Error writing to db: " + err.Error())
	}

	return c.Status(200).SendString("Success. Server saved to db")
}
