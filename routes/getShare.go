package routes

import (
	"openrsacloud/backend/db"
	"openrsacloud/backend/middlewares"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func GetShare(c *fiber.Ctx) error {
	shareId := c.Params("id")
	resp, err := db.DB.Select("shares:" + shareId)
	if err != nil {
		return err
	}
	var shareData db.Share
	err = surrealdb.Unmarshal(resp, &shareData)
	if err != nil {
		return err
	}

	if len(shareData.Recipients) != 0 {
		err = middlewares.NeedSession(c)
		if err != nil {
			return err
		}
		sessionData := c.Locals("session").(db.Session)
		if !slices.Contains(shareData.Recipients, sessionData.User) {
			return c.Status(401).JSON(fiber.Map{
				"status":  401,
				"message": "Unauthorized",
			})
		}
	}

	resp, err = db.DB.Select(shareData.Object)
	if err != nil {
		return err
	}
	var objectData fiber.Map
	err = surrealdb.Unmarshal(resp, &objectData)
	if err != nil {
		return err
	}

	if strings.HasPrefix(shareData.Object, "folder") {
		resp, err = db.DB.Query(`
			SELECT * FROM files WHERE parent = $folderId;
			SELECT * FROM folders WHERE parent = $folderId;
			`, fiber.Map{
			"folderId": shareData.Object,
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
				"parent":  objectData,
				"files":   folderItems[0].Result,
				"folders": folderItems[1].Result,
			},
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"status":  200,
			"message": "OK",
			"data":    objectData,
		})
	}
}
