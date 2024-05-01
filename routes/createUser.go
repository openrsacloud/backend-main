package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func CreateUser(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	var body map[string]interface{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
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

	if !userData.Admin {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "You are not an admin!",
		})
	}

	resp, err = db.DB.Query("CREATE users SET username = $username, password = crypto::argon2::generate($password), admin = $admin", fiber.Map{
		"username": body["username"],
		"password": body["password"],
		"admin":    body["admin"],
	})
	if err != nil {
		return err
	}
	var newUserData []surrealdb.RawQuery[[]db.User]
	err = surrealdb.Unmarshal(resp, &newUserData)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data":    newUserData[0].Result[0],
	})
}
