package routes

import (
	"log"
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetSessions(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	resp, err := db.DB.Query(
		`SELECT * FROM sessions WHERE user = $user`,
		map[string]interface{}{
			"user": sessionData.User,
		})
	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Internal server error",
			"info":    "Failed to get sessions",
		})
	}
	var allUserSessions []surrealdb.RawQuery[[]db.Session]
	err = surrealdb.Unmarshal(resp, &allUserSessions)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Internal server error",
			"info":    "Failed to parse sessions",
		})
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    allUserSessions[0].Result,
	})
}
