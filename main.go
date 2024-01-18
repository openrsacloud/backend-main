package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/surrealdb/surrealdb.go"
)

var JWTSecret []byte
var app *fiber.App
var db *surrealdb.DB

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
		log.Println(err.Error())
	}
	JWTSecret = []byte(os.Getenv("JWTSecret"))

	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "GoFiber",
		BodyLimit:    -1,
		AppName:      "OpenRSACloud_backend",
	})

	db, err := surrealdb.New("ws://truenas.local:30888/rpc")
	if err != nil {
		panic(err)
	}
	if _, err := db.Signin(map[string]interface{}{
		"ns":   "OpenRSACloud",
		"db":   "main",
		"user": os.Getenv("SurrealDatabaseUser"),
		"pass": os.Getenv("SurrealDatabasePass"),
	}); err != nil {
		panic(err)
	}
	if _, err := db.Use("OpenRSACloud", "main"); err != nil {
		panic(err)
	}

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Szia :3")
	})

	auth := api.Group("/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		var body map[string]string
		err := c.BodyParser(&body)
		if err != nil {
			c.Status(400)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Bad request",
				"info":    "Failed to parse request body",
			})
		}
		if body["username"] == "" || body["password"] == "" {
			c.Status(400)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Bad request",
				"info":    "The username or password is empty",
			})
		}
		resp, err := db.Query(`
			SELECT * FROM users WHERE username = $username 
			AND crypto::argon2::compare(password, $password);`,
			map[string]interface{}{
				"username": body["username"],
				"password": body["password"],
			})
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(fiber.Map{
				"status":  500,
				"message": "Internal server error",
				"info":    "Failed to get user data",
			})
		}
		var userData []surrealdb.RawQuery[[]map[string]interface{}]
		err = surrealdb.Unmarshal(resp, &userData)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(fiber.Map{
				"status":  500,
				"message": "Internal server error",
				"info":    "Failed to get user data",
			})
		}

		if len(userData) == 0 || len(userData[0].Result) == 0 {
			c.Status(401)
			return c.JSON(fiber.Map{
				"status":  401,
				"message": "Unauthorized",
				"info":    "The username or password is incorrect",
			})
		}

		delete(userData[0].Result[0], "password")

		resp, err = db.Query(`
			CREATE sessions SET
			user = $user,
			ip_address = $ip_address,
			end = time::now() + 2w,
			user_agent = $user_agent;`,
			map[string]interface{}{
				"user":       userData[0].Result[0]["id"],
				"user_agent": c.Get("User-Agent"),
				"ip_address": c.IP(),
			})
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(fiber.Map{
				"status":  500,
				"message": "Internal server error",
				"info":    "Failed to create session",
			})
		}
		var sessionData []surrealdb.RawQuery[[]map[string]interface{}]
		err = surrealdb.Unmarshal(resp, &sessionData)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(fiber.Map{
				"status":  500,
				"message": "Internal server error",
				"info":    "Failed to create session",
			})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"session": sessionData[0].Result[0]["id"],
		})
		signedToken, err := token.SignedString(JWTSecret)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(fiber.Map{
				"status":  500,
				"message": "Internal server error",
				"info":    "Failed to sign JsonWebToken",
			})
		}
		userData[0].Result[0]["token"] = signedToken

		c.Status(200)
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "OK",
			"data":    userData[0].Result[0],
		})
	})

	auth.Get("/getSessions", func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")[7:]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil {
			c.Status(401)
			return c.JSON(fiber.Map{
				"status":  401,
				"message": "Unauthorized",
				"info":    "Failed to parse JsonWebToken",
			})
		}
		tokenData := token.Claims.(jwt.MapClaims)
		resp, err := db.Select(tokenData["session"].(string))
		if err != nil {
			c.Status(401)
			return c.JSON(fiber.Map{
				"status":  401,
				"message": "Unauthorized",
				"info":    "Invalid session",
			})
		}
		var sessionData map[string]interface{}
		err = surrealdb.Unmarshal(resp, &sessionData)
		if err != nil {
			c.Status(401)
			return c.JSON(fiber.Map{
				"status":  401,
				"message": "Unauthorized",
				"info":    "Invalid session",
			})
		}
		resp, err = db.Query(
			`SELECT * FROM sessions WHERE user = $user`,
			map[string]interface{}{
				"user": sessionData["user"],
			})
		if err != nil {
			c.Status(500)
			return c.JSON(fiber.Map{
				"status":  500,
				"message": "Internal server error",
				"info":    "Failed to get sessions",
			})
		}
		var allUserSessions []surrealdb.RawQuery[[]map[string]interface{}]
		err = surrealdb.Unmarshal(resp, &allUserSessions)
		if err != nil {
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
	})

	app.Use(logger.New(logger.ConfigDefault))
	app.Listen("0.0.0.0:3250")
}
