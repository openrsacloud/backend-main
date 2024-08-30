package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetUser(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	resp, err := db.DB.Select(sessionData.User)
	if err != nil {
		return err
	}
	var userData db.User
	err = surrealdb.Unmarshal(resp, &userData)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    userData,
	})
}
