package routes

import (
	"openrsacloud/backend/db"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/surrealdb/surrealdb.go"
)

func Login(c *fiber.Ctx) error {
	var body map[string]string
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad request",
			"info":    "Failed to parse request body",
		})
	}
	if body["username"] == "" || body["password"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad request",
			"info":    "The username or password is empty",
		})
	}
	resp, err := db.DB.Query(`
		SELECT * FROM users WHERE username = $username 
		AND crypto::argon2::compare(password, $password);`,
		map[string]interface{}{
			"username": body["username"],
			"password": body["password"],
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
			"info":    "The username or password is incorrect",
		})
	}

	resp, err = db.DB.Query(`
		CREATE sessions SET
		user = $user,
		ip_address = $ip_address,
		end = time::now() + 2w,
		user_agent = $user_agent;`,
		map[string]interface{}{
			"user":       userData[0].Result[0].Id,
			"user_agent": c.Get("User-Agent"),
			"ip_address": c.IP(),
		})
	if err != nil {
		return err
	}
	var sessionData []surrealdb.RawQuery[[]db.Session]
	err = surrealdb.Unmarshal(resp, &sessionData)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session": sessionData[0].Result[0].Id,
	})
	signedToken, err := token.SignedString([]byte(os.Getenv("JWTSecret")))
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
		"data": fiber.Map{
			"user":  userData[0].Result[0],
			"token": signedToken,
		},
	})
}
