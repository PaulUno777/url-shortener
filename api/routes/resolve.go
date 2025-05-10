package routes

import (
	"github.com/PaulUno777/url-shortener/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	db := database.CreateClient(0)
	defer db.Close()

	value, err := db.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something when wrong",
		})
	}

	db2 := database.CreateClient(1)
	defer db.Close()

	db2.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
