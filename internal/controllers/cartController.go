package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"github.com/redis/go-redis/v9"
)

type CartController interface {

	// Get cart by Id
	GetCart(c *fiber.Ctx) error

	// Add selected product(s) to cart
	UpdateCart(c *fiber.Ctx) error
}

var (
	cartInstance *cartController
	cart         models.Cart
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

func (cc *cartController) GetCart(c *fiber.Ctx) error {
	accountId := c.Query("id")
	if accountId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request",
		})
	}

	if err := cc.db.UseGorm().First(&cart, accountId); err != nil {
		return c.SendString("error: Unable to find account")
	}

	result, _ := json.Marshal(&cart)

	return c.SendString(string(result))

}

func (cc *cartController) UpdateCart(c *fiber.Ctx) error {
	requestBody := new(cartRequest)
	if err := c.BodyParser(&requestBody); err != nil {
		return err
	}

	return cc.db.UseGorm().Save(&requestBody).Error
}
