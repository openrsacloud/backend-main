package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func RemoveUser(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	var body struct {
		UserId string `json:"user_id"`
	}
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad request",
			"info":    "Failed to parse request body",
		})
	}

	resp, err := db.DB.Select(sessionData.User)
	if err != nil {
		return err
	}
	var userData db.User
	err = surrealdb.Unmarshal(resp, userData)
	if err != nil {
		return err
	}

	if !userData.Admin || body.UserId == sessionData.User {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "You cannot remove this user!",
		})
	}

	_, err = db.DB.Delete(body.UserId)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})
}
