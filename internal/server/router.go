package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kevinhartarto/market-be/internal/controllers"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/redis/go-redis/v9"
)

func NewHandler(db database.Service, redis *redis.Client) *fiber.App {
	app := fiber.New()

	app.Use(healthcheck.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	marketAPI := app.Group("/api")

	// Login and Register APIs
	identity := controllers.NewIdentityController(db, redis)

	identityAPI := marketAPI.Group("/user")
	identityAPI.Post("/login", func(c *fiber.Ctx) error {
		return identity.Login(c)
	})
	identityAPI.Post("/register", func(c *fiber.Ctx) error {
		return identity.CreateUser(c)
	})
	identityAPI.Delete("/delete", func(c *fiber.Ctx) error {
		return identity.Deactivate(c)
	})

	// For Roles higher than viewer
	superUserAPI := identityAPI.Group("/super")
	superUserAPI.Post("/register", func(c *fiber.Ctx) error {
		return identity.CreateAdmin(c)
	})
	superUserAPI.Post("/role", func(c *fiber.Ctx) error {
		return identity.CreateRole(c)
	})

	return app
}
