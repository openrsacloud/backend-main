package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func RemoveUser(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	var body struct {
		Password string `json:"password"`
	}
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad request",
			"info":    "Failed to parse request body",
		})
	}

	resp, err := db.DB.Query(`
		SELECT * FROM users WHERE id = $uid 
		AND crypto::argon2::compare(password, $password);`,
		fiber.Map{
			"uid":      sessionData.User,
			"password": body.Password,
		})
	if err != nil {
		return err
	}
	var userData []surrealdb.RawQuery[[]db.User]
	err = surrealdb.Unmarshal(resp, &userData)
	if err != nil {
		return err
	}

	if len(userData) == 0 || len(userData[0].Result) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "The password is incorrect",
		})
	}

	_, err = db.DB.Delete(sessionData.User)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})
}
