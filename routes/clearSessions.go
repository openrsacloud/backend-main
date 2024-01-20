package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
)

func ClearSessions(c *fiber.Ctx) error {
	sessionData := c.Locals("session")
	if sessionData == nil {
		return nil
	}
	_, err := db.DB.Query("DELETE sessions WHERE user = $user", fiber.Map{
		"user": sessionData.(db.Session).User,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": "Internal server error",
			"info":    "Failed to delete sessions",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})
}
