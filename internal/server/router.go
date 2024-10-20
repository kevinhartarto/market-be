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
		return account.GetAccount(c)
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

	// For admin and owner roles
	superUserAPI := accountAPI.Group("/super")
	superUserAPI.Get("/reload", func(c *fiber.Ctx) error {
		account.GetAllRoles(context)
		return c.SendStatus(fiber.StatusOK)
	})
	superUserAPI.Post("/role", func(c *fiber.Ctx) error {
		return account.CreateRole(c)
	})
	superUserAPI.Put("/update", func(c *fiber.Ctx) error {
		return account.UpdateRole(c)
	})

	// Cart
	cart := controllers.NewCartController(db, redis)
	cartAPI := marketAPI.Group("/cart")
	cartAPI.Get("/", func(c *fiber.Ctx) error {
		return cart.GetCart(c)
	})
	cartAPI.Put("/update", func(c *fiber.Ctx) error {
		return cart.UpdateCart(c)
	})

	product := controllers.NewProductController(db, redis)
	productAPI := marketAPI.Group("/product")
	productAPI.Get("/", func(c *fiber.Ctx) error {
		return product.GetAllProducts(c)
	})
	productAPI.Get("/brands", func(c *fiber.Ctx) error {
		return product.GetAllBrands(c)
	})
	productAPI.Get("/categories", func(c *fiber.Ctx) error {
		return product.GetAllCategories(c)
	})

	productAPI.Get("/brand/products", func(c *fiber.Ctx) error {
		return product.GetProductsByBrand(c)
	})
	productAPI.Get("/category/products", func(c *fiber.Ctx) error {
		return product.GetProductsByCategory(c)
	})

	productAPI.Get("/detail", func(c *fiber.Ctx) error {
		return product.GetProductDetails(c)
	})
	productAPI.Get("/brand/detail", func(c *fiber.Ctx) error {
		return product.GetBrandDetails(c)
	})
	productAPI.Get("/category/detail", func(c *fiber.Ctx) error {
		return product.GetCategoryDetails(c)
	})

	productAPI.Put("/update", func(c *fiber.Ctx) error {
		return product.UpdateProduct(c)
	})
	productAPI.Put("/brand/update", func(c *fiber.Ctx) error {
		return product.UpdateBrand(c)
	})
	productAPI.Put("/category/update", func(c *fiber.Ctx) error {
		return product.UpdateCategory(c)
	})

	return app
}
