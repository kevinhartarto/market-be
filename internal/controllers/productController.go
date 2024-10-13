package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"github.com/redis/go-redis/v9"
)

type ProductController interface {

	// Retrieve all products
	GetAllProducts(c *fiber.Ctx) error
}

var (
	productInstance *productController
)

type productController struct {
	db    database.Service
	redis *redis.Client
}

func NewProductController(db database.Service, redis *redis.Client) *productController {

	if productInstance != nil {
		return productInstance
	}

	productInstance = &productController{
		db:    db,
		redis: redis,
	}

	return productInstance
}

func (pc *productController) GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product
	pc.db.UseGorm().Where("active").Find(&products)

	result, _ := json.Marshal(products)

	return c.SendString(string(result))
}
