package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/redis/go-redis/v9"
)

type CartController interface {

	// Add selected product(s) to cart
	UpdateCart(c *fiber.Ctx) error
}

var (
	cartInstance *cartController
)

type cartController struct {
	db    database.Service
	redis *redis.Client
}

type cartRequest struct {
	Id       uuid.UUID
	Products []string
}

func NewCartController(db database.Service, redis *redis.Client) *cartController {

	if cartInstance != nil {
		return cartInstance
	}

	cartInstance = &cartController{
		db:    db,
		redis: redis,
	}

	return cartInstance
}

func (cc *cartController) UpdateCart(c *fiber.Ctx) error {
	requestBody := new(cartRequest)
	if err := c.BodyParser(&requestBody); err != nil {
		return err
	}

	return cc.db.UseGorm().Save(&requestBody).Error
}
