package main

import (
	"errors"
	"log"
	"os"

	"openrsacloud/backend/db"
	"openrsacloud/backend/middlewares"
	"openrsacloud/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
		log.Println(err.Error())
	}

	if _, ok := os.LookupEnv("BasePath"); !ok {
		panic("BasePath enviornmnet variable is not set")
	}
	if _, ok := os.LookupEnv("JWTSecret"); !ok {
		panic("JWTSecret enviornmnet variable is not set")
	}

	db.Connect()

	app := fiber.New(fiber.Config{
		Prefork:           true,
		BodyLimit:         1099511627776,
		StreamRequestBody: true,
		AppName:           "OpenRSACloud_backend",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				return c.Status(e.Code).JSON(fiber.Map{
					"status":  e.Code,
					"message": e.Message,
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"status":  500,
				"message": "Internal Server Error",
				"info":    err.Error(),
			})
		},
	})
	app.Use(cors.New())
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

	defer app.ShutdownWithTimeout(2000)
	defer db.Disconnect()
}

func initRoutes(r fiber.Router) {
	auth := r.Group("/auth")
	auth.Post("/login", routes.Login)
	auth.Get("/create_user", middlewares.NeedSession, routes.CreateUser)
	auth.Get("/get_user", middlewares.NeedSession, routes.GetAccount)
	auth.Get("/get_sessions", middlewares.NeedSession, routes.GetSessions)
	auth.Post("/clear_sessions", middlewares.NeedSession, routes.ClearSessions)
	auth.Post("/remove_session", middlewares.NeedSession, routes.RemoveSession)

	files := r.Group("/files")
	files.Post("/upload", middlewares.NeedSession, routes.UploadFile)
	files.Get("/get_file/:id", routes.GetFile)
	files.Post("/create_folder", middlewares.NeedSession, routes.CreateFolder)
	files.Get("/get_folder/:id?", middlewares.NeedSession, routes.GetFolder)

	shares := r.Group("/share")
	shares.Post("/create_share", middlewares.NeedSession, routes.CreateShare)
	shares.Post("/remove_share", middlewares.NeedSession, routes.RemoveShare)
	shares.Post("/update_share", middlewares.NeedSession, routes.UpdateShare)
	shares.Get("/:id", routes.GetShare)
}
