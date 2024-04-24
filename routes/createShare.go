package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func CreateShare(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	var body map[string]string
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	resp, err := db.DB.Create("shares", fiber.Map{
		"object":     body["object"],
		"owner":      sessionData.User,
		"recipients": body["recipients"],
	})
	if err != nil {
		return err
	}
	var shareData db.Share
	err = surrealdb.Unmarshal(resp, &shareData)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    shareData,
	})
}
