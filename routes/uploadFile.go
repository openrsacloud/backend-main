package routes

import (
	"io"
	"openrsacloud/backend/db"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func UploadFile(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request",
			"info":    "Failed to parse form data",
		})
	}

	_newFileOptions := fiber.Map{
		"name":  c.FormValue("filename", file.Filename),
		"owner": sessionData.User,
		"size":  file.Size,
		"type":  file.Header.Get("Content-Type"),
	}
	if c.FormValue("parent") != "" {
		_newFileOptions["parent"] = c.FormValue("parent")
	}
	resp, err := db.DB.Create("files", _newFileOptions)
	if err != nil {
		return err
	}

	var fileData []db.File
	err = surrealdb.Unmarshal(resp, &fileData)
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request",
			"info":    "Failed to open file",
		})
	}
	defer src.Close()

	destination, err := os.Create(os.Getenv("BasePath") + fileData[0].Id[7:])
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, src)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})
}
