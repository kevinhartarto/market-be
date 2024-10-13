package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kevinhartarto/market-be/internal/controllers"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/redis/go-redis/v9"
)

func NewHandler(db database.Service, redis *redis.Client) *fiber.App {
	context := context.Background()
	app := fiber.New()

	app.Use(healthcheck.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	marketAPI := app.Group("/api")

	// Login and Register APIs
	account := controllers.NewAccountController(db, redis)
	account.GetAllRoles(context)
	fmt.Println("Roles loaded")

	accountAPI := marketAPI.Group("/user")
	accountAPI.Get("/", func(c *fiber.Ctx) error {
		return account.ShowAccountDetails(c)
	})
	accountAPI.Post("/login", func(c *fiber.Ctx) error {
		return account.Login(c)
	})
	accountAPI.Post("/register", func(c *fiber.Ctx) error {
		return account.CreateAccount(c)
	})
	accountAPI.Put("/update", func(c *fiber.Ctx) error {
		return account.UpdateAccount(c)
	})
	accountAPI.Put("/verify", func(c *fiber.Ctx) error {
		return account.ChangeRole(c)
	})
	accountAPI.Delete("/delete", func(c *fiber.Ctx) error {
		return account.Deactivate(c)
	})

	// For admin and owner roles
	superUserAPI := accountAPI.Group("/super")
	superUserAPI.Get("/reload", func(c *fiber.Ctx) error {
		account.GetAllRoles(context)
		return c.SendStatus(fiber.StatusOK)
	})
	superUserAPI.Post("/register", func(c *fiber.Ctx) error {
		return account.CreateAdmin(c)
	})
	superUserAPI.Post("/role", func(c *fiber.Ctx) error {
		return account.CreateRole(c)
	})
	superUserAPI.Put("/permissions", func(c *fiber.Ctx) error {
		return account.ChangeRolePermissions(c)
	})
	superUserAPI.Put("/elevate", func(c *fiber.Ctx) error {
		return account.SetRoleToAdmin(c)
	})

	// Cart
	cart := controllers.NewCartController(db, redis)
	marketAPI.Put("/cart", func(c *fiber.Ctx) error {
		return cart.UpdateCart(c)
	})

	product := controllers.NewProductController(db, redis)
	productAPI := marketAPI.Group("/product")
	productAPI.Get("/", func(c *fiber.Ctx) error {
		return product.GetAllProducts(c)
	})

	return app
}
