package routes

import (
	"log"
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetAccount(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	resp, err := db.DB.Select(sessionData.User)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": "Internal server error",
			"info":    "Failed to get user",
		})
	}
	var userData db.User
	err = surrealdb.Unmarshal(resp, &userData)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": "Internal server error",
			"info":    "Failed to get user",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    userData,
	})
}
