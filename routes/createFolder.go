package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func CreateFolder(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	var body map[string]interface{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	resp, err := db.DB.Create("folders", fiber.Map{
		"name":   body["name"],
		"parent": body["parent"],
		"owner":  sessionData.User,
	})
	if err != nil {
		return err
	}
	var folderData db.Folder
	err = surrealdb.Unmarshal(resp, folderData)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    folderData,
	})
}
