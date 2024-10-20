package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"github.com/redis/go-redis/v9"
)

type ProductController interface {

	// Retrieve all products
	GetAllBrands(c *fiber.Ctx) error

	GetBrandDetails(c *fiber.Ctx) error

	UpdateBrand(c *fiber.Ctx) error

	GetAllCategories(c *fiber.Ctx) error

	GetCategoryDetails(c *fiber.Ctx) error

	UpdateCategories(c *fiber.Ctx) error

	GetAllProducts(c *fiber.Ctx) error

	GetProductsByBrand(c *fiber.Ctx) error

	GetProductsByCategory(c *fiber.Ctx) error

	GetProductDetails(c *fiber.Ctx) error

	UpdateProduct(c *fiber.Ctx) error
}

var (
	productInstance *productController

	brands     []models.Brand
	categories []models.Category
	products   []models.Product

	brand    models.Brand
	category models.Category
	product  models.Product
)

type productController struct {
	db    database.Service
	redis *redis.Client
}

type BrandToUpdate struct {
	Id        uuid.UUID
	NewStatus bool
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
	pc.db.UseGorm().Where("active").Find(&brands)
	result, _ := json.Marshal(brands)
	return c.SendString(string(result))
}

func (pc *productController) GetBrandDetails(c *fiber.Ctx) error {
	if err := c.BodyParser(&brand); err != nil {
		return err
	}

	if err := pc.db.UseGorm().First(&brand).Error; err != nil {
		return err
	}

	result, _ := json.Marshal(&brand)
	return c.SendString(string(result))
}

func (pc *productController) UpdateBrand(c *fiber.Ctx) error {
	var updateBrand struct {
		brand       models.Brand
		updateType  string
		updateValue bool
	}
	success := false

	if err := c.BodyParser(&updateBrand); err != nil {
		return err
	}

	switch updateBrand.updateType {
	case "update":
		affectedRows = pc.db.UseGorm().Save(&updateBrand.brand).RowsAffected
	case "sale":
		affectedRows = pc.db.UseGorm().Model(&updateBrand.brand).Update("on_sale", updateBrand.updateValue).RowsAffected
	case "active":
		affectedRows = pc.db.UseGorm().Model(&updateBrand.brand).Update("active", updateBrand.updateValue).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&updateBrand.brand)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (pc *productController) GetAllCategories(c *fiber.Ctx) error {
	pc.db.UseGorm().Where("active").Find(&categories)
	result, _ := json.Marshal(categories)
	return c.SendString(string(result))
}

func (pc *productController) GetCategoryDetails(c *fiber.Ctx) error {
	if err := c.BodyParser(&category); err != nil {
		return err
	}

	if err := pc.db.UseGorm().First(&category).Error; err != nil {
		return err
	}

	result, _ := json.Marshal(&category)
	return c.SendString(string(result))
}

func (pc *productController) UpdateCategory(c *fiber.Ctx) error {
	var updateCategory struct {
		category    models.Category
		updateType  string
		updateValue bool
	}
	success := false

	if err := c.BodyParser(&updateCategory); err != nil {
		return err
	}

	switch updateCategory.updateType {
	case "update":
		affectedRows = pc.db.UseGorm().Save(&updateCategory.category).RowsAffected
	case "featured":
		affectedRows = pc.db.UseGorm().Model(&updateCategory.category).Update("featured", updateCategory.updateValue).RowsAffected
	case "active":
		affectedRows = pc.db.UseGorm().Model(&updateCategory.category).Update("active", updateCategory.updateValue).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&updateCategory.category)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (pc *productController) GetAllProducts(c *fiber.Ctx) error {
	pc.db.UseGorm().Where("active").Find(&products)
	result, _ := json.Marshal(products)
	return c.SendString(string(result))
}

func (pc *productController) GetProductsByBrand(c *fiber.Ctx) error {
	brand := c.Query("brand")
	pc.db.UseGorm().Where("brand = ?", brand).Find(&products)
	result, _ := json.Marshal(products)
	return c.SendString(string(result))
}

func (pc *productController) GetProductsByCategory(c *fiber.Ctx) error {
	category := c.Query("category")
	pc.db.UseGorm().Where("category = ?", category).Find(&products)
	result, _ := json.Marshal(products)
	return c.SendString(string(result))
}

func (pc *productController) GetProductDetails(c *fiber.Ctx) error {
	if err := c.BodyParser(&product); err != nil {
		return err
	}

	if err := pc.db.UseGorm().First(&product).Error; err != nil {
		return err
	}

	result, _ := json.Marshal(&product)
	return c.SendString(string(result))
}

func (pc *productController) UpdateProduct(c *fiber.Ctx) error {
	var updateProduct struct {
		product     models.Product
		updateType  string
		updateValue bool
	}
	success := false

	if err := c.BodyParser(&updateProduct); err != nil {
		return err
	}

	switch updateProduct.updateType {
	case "update":
		affectedRows = pc.db.UseGorm().Save(&updateProduct.product).RowsAffected
	case "active":
		affectedRows = pc.db.UseGorm().Model(&updateProduct.product).Update("active", updateProduct.updateValue).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&updateProduct.product)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}
