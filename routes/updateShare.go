package routes

import (
	"openrsacloud/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

func UpdateShare(c *fiber.Ctx) error {
	sessionData := c.Locals("session").(db.Session)
	var body map[string]string
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	resp, err := db.DB.Query("UPDATE $share SET recipients = $recipients WHERE owner = $user", fiber.Map{
		"user":       sessionData.User,
		"share":      body["id"],
		"recipients": body["recipients"],
	})
	if err != nil {
		return err
	}
	var shareData []surrealdb.RawQuery[[]db.Share]
	err = surrealdb.Unmarshal(resp, &shareData)
	if err != nil {
		return err
	}

	if len(shareData[0].Result) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  404,
			"message": "Not found",
			"info":    "A share with the porvided id was not found.",
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"status":  200,
			"message": "OK",
			"data":    shareData[0].Result[0],
		})
	}

}
