package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/PaulUno777/url-shortener/database"
	"github.com/PaulUno777/url-shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateLimitRest time.Duration `json:"rate_limit_rest"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	println("Managing rate limit")
	//Manage rate limit
	db2 := database.CreateClient(1)
	defer db2.Close()
	value, err := db2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		db2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to database",
		})
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something when wrong while converting",
		})
	}
	if valueInt <= 0 {
		limit, _ := db2.TTL(database.Ctx, c.IP()).Result()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":           "Rate limit exceeded",
			"rate_limit_rest": limit / time.Nanosecond / time.Minute,
		})
	}

	//check if the input is URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	//check domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL DOMAIN",
		})
	}

	//enforce https SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	//shorten URL logic
	var id string
	if body.CustomShort == "" {
		id = uuid.NewString()[:6]
	} else {
		id = body.CustomShort
	}

	db1 := database.CreateClient(0)
	defer db1.Close()

	value, _ = db1.Get(database.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short already used",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	err = db1.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to database",
		})
	}

	res := response{
		URL:            body.URL,
		CustomShort:    "",
		Expiry:         body.Expiry,
		XRateRemaining: 10,
		XRateLimitRest: 30,
	}

	db2.Decr(database.Ctx, c.IP())

	value, _ = db2.Get(database.Ctx, c.IP()).Result()
	res.XRateRemaining, _ = strconv.Atoi(value)

	ttl, _ := db2.TTL(database.Ctx, c.IP()).Result()
	res.XRateLimitRest = ttl / time.Nanosecond / time.Minute

	res.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(res)

}
