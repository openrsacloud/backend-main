package routes

import (
	"openrsacloud/backend/db"
	"openrsacloud/backend/middlewares"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetFile(c *fiber.Ctx) error {
	fileId := c.Params("id")
	if !strings.HasPrefix(fileId, "files:") {
		return c.Status(401).JSON(fiber.Map{
			"status":  400,
			"message": "Bad request",
			"info":    "Invalid file ID",
		})
	}
	resp, err := db.DB.Select(fileId)
	if err != nil {
		return err
	}
	var fileData db.File
	err = surrealdb.Unmarshal(resp, &fileData)
	if err != nil {
		return err
	}

	resp, err = db.DB.Query("SELECT * FROM shares WHERE object = $obj", fiber.Map{
		"obj": "files:" + fileId,
	})
	if err != nil {
		return err
	}
	var shareData []surrealdb.RawQuery[[]db.Share]
	err = surrealdb.Unmarshal(resp, &fileData)
	if err != nil {
		return err
	}
	if len(shareData[0].Result) == 0 || len(shareData[0].Result[0].Recipients) != 0 {
		err = middlewares.NeedSession(c)
		if err != nil {
			return err
		}
		sessionData := c.Locals("session").(db.Session)
		if sessionData.User != fileData.Owner && !slices.Contains(shareData[0].Result[0].Recipients, sessionData.User) {
			return c.Status(401).JSON(fiber.Map{
				"status":  401,
				"message": "Unauthorized",
			})
		}
	}

	return c.Status(200).Download(os.Getenv("BasePath")+fileId, fileData.Name)
}
