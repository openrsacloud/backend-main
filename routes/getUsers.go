package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetUsers(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)

	resp, err := db.DB.Select(sessionData.User)
	if err != nil {
		return err
	}
	var sessionUser db.User
	err = surrealdb.Unmarshal(resp, sessionUser)
	if err != nil {
		return err
	}

	if !sessionUser.Admin {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "You are not an admin!",
		})
	}

	resp, err = db.DB.Query("SELECT * FROM users LIMIT 50 START $start", fiber.Map{
		"start": c.QueryInt("page", 0)*50 + 1,
	})
	if err != nil {
		return err
	}
	var userData []surrealdb.RawQuery[[]db.User]
	err = surrealdb.Unmarshal(resp, userData)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    userData[0].Result,
	})
}
