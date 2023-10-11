package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "OpenRSACloud",
		BodyLimit:    -1,
		AppName:      "OpenRSACloud backend",
	})

	app.Use(logger.New(logger.ConfigDefault))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Szia:)")
	})

	app.Listen("0.0.0.0:3250")
}
