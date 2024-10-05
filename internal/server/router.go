package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"gorm.io/gorm/logger"
)

func NewHandler(db database.Service) *fiber.App {
	app := fiber.New()

	app.Use(healthcheck.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	marketAPI := app.Group("/api")

	// Login and Register APIs
	identity := models.CreateIdentity(db)

	identityAPI := marketAPI.Group("/identity")
	identityAPI.Post("/auth", func(c *fiber.Ctx) error {
		return identity.Auth()
	})
	identityAPI.Post("/login", func(c *fiber.Ctx) error {
		return identity.Login()
	})
	identityAPI.Post("/register", func(c *fiber.Ctx) error {
		return identity.Register()
	})

	return app
}
