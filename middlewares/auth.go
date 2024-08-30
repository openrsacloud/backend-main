package middlewares

import (
	"openrsacloud/backend/db"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/surrealdb/surrealdb.go"
)

func NeedSession(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if len(tokenString) < 10 {
		c.Locals("session", nil)
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "Invalid token",
		})
	}
	token, err := jwt.Parse(tokenString[7:], func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTSecret")), nil
	})
	if err != nil {
		c.Locals("session", nil)
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "Invalid token",
		})
	}
	tokenClaims := token.Claims.(jwt.MapClaims)
	resp, err := db.DB.Select(tokenClaims["session"].(string))
	if err != nil {
		c.Locals("session", nil)
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "Invalid token",
		})
	}
	var sessionData db.Session
	err = surrealdb.Unmarshal(resp, &sessionData)
	if err != nil {
		c.Locals("session", nil)
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
			"info":    "Invalid token",
		})
	}
	c.Locals("session", sessionData)
	return c.Next()
}
