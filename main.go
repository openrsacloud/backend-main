package main

import (
	"log"

	"openrsacloud/backend/db"
	"openrsacloud/backend/middlewares"
	"openrsacloud/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
		log.Println(err.Error())
	}

	db.Connect()
	defer db.Disconnect()

	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "GoFiber",
		BodyLimit:    -1,
		AppName:      "OpenRSACloud_backend",
	})
	app.Use(logger.New(logger.ConfigDefault))

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		c.Status(200)
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "OK",
			"info":    "You have reached the OpenRSACloud API c:",
		})
	})

	initRoutes(api)

	app.Listen("0.0.0.0:3250")
}

func initRoutes(r fiber.Router) {
	auth := r.Group("/auth")
	auth.Post("/login", routes.Login)
	auth.Get("/get_sessions", middlewares.NeedSession, routes.GetSessions)
	auth.Get("/clear_sessions", middlewares.NeedSession, routes.ClearSessions)
}
