package main

import (
	"encoding/json"
	"os"

	"github.com/alexander-grube/secret/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/secret", func(c *fiber.Ctx) error {
		// use json from context body
		secret := &model.Secret{}
		id := uuid.New().String()
		secret.ID = id
		if err := c.BodyParser(secret); err != nil {
			return err
		}

		s, err := json.Marshal(secret)
		if err != nil {
			return err
		}

		if err := rdb.Set(c.Context(), id, s, secret.TTL).Err(); err != nil {
			return err
		}

		return c.SendString(id)
	})

	app.Get("/secret/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		secret := &model.Secret{}

		if secretJson, err := rdb.Get(c.Context(), id).Result(); err == nil {
			if err = rdb.Del(c.Context(), id).Err(); err != nil {
				return err
			}
			json.Unmarshal([]byte(secretJson), &secret)
			return c.SendString(secret.Data)
		}

		return c.Status(fiber.StatusNotFound).SendString("Secret not found")
	})

	app.Listen(PORT)
}
