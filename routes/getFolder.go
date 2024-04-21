package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetFolder(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	folderId := c.Params("id")

	if folderId != "" {
		resp, err := db.DB.Select(folderId)
		if err != nil {
			return err
		}

		var folderData db.Folder
		err = surrealdb.Unmarshal(resp, &folderData)
		if err != nil {
			return err
		}

		if folderData.Owner != sessionData.User {
			return fiber.ErrNotFound
		}

		resp, err = db.DB.Query(`
			SELECT * FROM files WHERE parent = $folderId AND owner = $user;
			SELECT * FROM folders WHERE parent = $folderId AND owner = $user;
			`, fiber.Map{
			"folderId": folderId,
			"user":     sessionData.User,
		})
		if err != nil {
			return err
		}

		var folderItems []surrealdb.RawQuery[[]interface{}]
		err = surrealdb.Unmarshal(resp, &folderItems)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  200,
			"message": "OK",
			"data": fiber.Map{
				"parent":  folderData,
				"files":   folderItems[0].Result,
				"folders": folderItems[1].Result,
			},
		})
	} else {

		resp, err := db.DB.Query(`
		SELECT * FROM files WHERE parent = NONE AND owner = $user;
		SELECT * FROM folders WHERE parent = NONE AND owner = $user;
		`, fiber.Map{
			"user": sessionData.User,
		})
		if err != nil {
			return err
		}

		var folderItems []surrealdb.RawQuery[[]interface{}]
		err = surrealdb.Unmarshal(resp, &folderItems)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  200,
			"message": "OK",
			"data": fiber.Map{
				"files":   folderItems[0].Result,
				"folders": folderItems[1].Result,
			},
		})
	}
}
