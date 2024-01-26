package routes

import (
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
		return err
	}
	var allUserSessions []surrealdb.RawQuery[[]db.Session]
	err = surrealdb.Unmarshal(resp, &allUserSessions)
	if err != nil {
		return err
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    allUserSessions[0].Result,
	})
}
