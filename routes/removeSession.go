package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
)

func RemoveSession(c *fiber.Ctx) error {
	_ = c.Locals("session").(db.Session)
	var body map[string]string
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad request",
			"info":    "Failed to parse request body",
		})
	}
	_, err = db.DB.Delete(body["session_id"])
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})
}
