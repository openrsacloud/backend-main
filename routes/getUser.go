package routes

import (
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	sessionData := c.Locals("session")
	if sessionData == nil {
		return nil
	}
	return c.Status(501).JSON(fiber.Map{
		"status":  501,
		"message": "Not implemented",
	})
}
