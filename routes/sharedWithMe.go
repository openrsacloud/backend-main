package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func SharedWithMe(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)

	resp, err := db.DB.Query("SELECT * FROM shares WHERE $user IN recipients;", fiber.Map{
		"user": sessionData.User,
	})
	if err != nil {
		return err
	}
	var sharesData []surrealdb.RawQuery[[]db.Share]
	err = surrealdb.Unmarshal(resp, &sharesData)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    sharesData[0].Result,
	})

}
