package main

import (
	"os"
	"time"

	"github.com/alexander-grube/secret/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/go-redis/redis/v8"
)

var (
	PORT string = ":" + os.Getenv("PORT")
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New())

	rdbOptions, err := redis.ParseURL(os.Getenv("REDIS_URL"))

	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(rdbOptions)

	secret := &model.Secret{
		ID:     "1",
		Data:   "secret",
		TTL:    time.Hour,
		IsSeen: false,
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/secret", func(c *fiber.Ctx) error {
		return rdb.Set(c.Context(), secret.ID, secret, secret.TTL).Err()
	})

	app.Get("/secret/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if secret, err := rdb.Get(c.Context(), id).Result(); err == nil {
			return c.JSON(secret)
		}

		return c.Status(fiber.StatusNotFound).SendString("Secret not found")
	})

	app.Listen(PORT)
}
