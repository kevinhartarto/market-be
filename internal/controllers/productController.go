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

func (pc *productController) GetAllBrands(c *fiber.Ctx) error {
	var brands []models.Brand
	pc.db.UseGorm().Where("active").Find(&brands)

	result, _ := json.Marshal(brands)

	return c.SendString(string(result))
}

func (pc *productController) GetBrandDetails(c *fiber.Ctx) error {
	brand := new(models.Brand)
	if err := c.BodyParser(&brand); err != nil {
		return err
	}

	return pc.db.UseGorm().First(&brand).Error
}

func (pc *productController) UpdateBrand(c *fiber.Ctx) error {
	brand := new(models.Brand)
	if err := c.BodyParser(&brand); err != nil {
		return err
	}

	return pc.db.UseGorm().Save(&brand).Error
}

func (pc *productController) SetBrandOnSale(c *fiber.Ctx) error {

}

func (pc *productController) SetBrandOffSale(c *fiber.Ctx) error {

}

func (pc *productController) DeactiveBrand(c *fiber.Ctx) error {

}

func (pc *productController) ReactiveBrand(c *fiber.Ctx) error {

}

func (pc *productController) GetAllCategories(c *fiber.Ctx) error {
	var categories []models.Category
	pc.db.UseGorm().Where("active").Find(&categories)

	result, _ := json.Marshal(categories)

	return c.SendString(string(result))
}

func (pc *productController) GetCategoryDetails(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(&category); err != nil {
		return err
	}

	return pc.db.UseGorm().First(&category).Error
}

func (pc *productController) UpdateCategory(c *fiber.Ctx) error {
	brand := new(models.Brand)
	if err := c.BodyParser(&brand); err != nil {
		return err
	}

	return pc.db.UseGorm().Save(&brand).Error
}

func (pc *productController) SetCategoryOnFeatured(c *fiber.Ctx) error {

}

func (pc *productController) SetCategoryOffFeatured(c *fiber.Ctx) error {

}

func (pc *productController) DeactiveCategory(c *fiber.Ctx) error {

}

func (pc *productController) ReactiveCategory(c *fiber.Ctx) error {

}

func (pc *productController) GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product
	pc.db.UseGorm().Where("active").Find(&products)

	result, _ := json.Marshal(products)

	return c.SendString(string(result))
}

func (pc *productController) GetProductDetails(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(&product); err != nil {
		return err
	}

	return pc.db.UseGorm().First(&product).Error
}

func (pc *productController) UpdateProduct(c *fiber.Ctx) error {

}

func (pc *productController) ChangeProductBrand(c *fiber.Ctx) error {

}

func (pc *productController) ChangeProductCategory(c *fiber.Ctx) error {

}

func (pc *productController) SetProductOnSale(c *fiber.Ctx) error {

}

func (pc *productController) SetProductOffSale(c *fiber.Ctx) error {

}

func (pc *productController) DeactiveProduct(c *fiber.Ctx) error {

}

func (pc *productController) ReactiveProduct(c *fiber.Ctx) error {

}
