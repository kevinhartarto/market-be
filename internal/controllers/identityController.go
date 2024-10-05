package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"github.com/kevinhartarto/market-be/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Identity Controller represent a service that handle all things related to Identity
type IdentityController interface {

	// Create an user identity.
	// It returns an error if the identity cannot be created or already exists.
	CreateUser(*fiber.Ctx, database.Service) error

	// Create an admin identity.
	// It returns an error if the identity cannot be created or already exists.
	CreateAdmin(*fiber.Ctx, database.Service) error

	// Log in an identity.
	// It returns an error if the identity cannot be identified.
	Login(*fiber.Ctx, database.Service) error

	// Create an identity role.
	// It returns an error if similar role exits
	CreateRole(*fiber.Ctx, database.Service) error

	// Deactivate an identity.
	// It returns an error if the identity not found.
	Deactivate(*fiber.Ctx, database.Service) error
}

type identityController struct {
	db    gorm.DB
	redis redis.Client
}

type tempUser struct {
	email    string `json:"email"`
	password string `json:"password"`
}

var (
	controllerInstance *identityController
)

func NewIdentityController(db gorm.DB, redis redis.Client) IdentityController {

	if controllerInstance != nil {
		return controllerInstance
	}

	controllerInstance = &identityController{
		db:    db,
		redis: redis,
	}

	return controllerInstance
}

func (ic *identityController) CreateUser(c *fiber.Ctx, db database.Service) error {

	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	// Hash the password for the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Server Error")
	}

	user.Id = uuid.New()
	user.Password = hashedPassword
	user.Role = getRole(db, "user")
	user.Verified = false
	user.Active = true

	return db.UseGorm().Create(&user).Error
}

func (ic *identityController) CreateAdmin(c *fiber.Ctx, db database.Service) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	user.Id = uuid.New()
	user.Role = getRole(db, "admin")
	user.Verified = false
	user.Active = true

	return db.UseGorm().Create(&user).Error
}

func (ic *identityController) Login(c *fiber.Ctx, db database.Service) error {

	tempUser := new(tempUser)
	user := new(models.User)

	if err := c.BodyParser(&tempUser); err != nil {
		return err
	}

	if err := db.UseGorm().Where("email = ? and active", user.Id).First(&user).Error; err != nil {
		return err
	}

	if err := utils.VerifyPassword(tempUser.password, user.Password).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	return c.Redirect("/")
}

func (ic *identityController) Deactivate(c *fiber.Ctx, db database.Service) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	if err := db.UseGorm().First(&role, role.Id).Error; err != nil {
		return err
	}

	return db.UseGorm().Model(&role).Update("Active", false).Error
}

func (ic *identityController) CreateRole(c *fiber.Ctx, db database.Service) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	return db.UseGorm().Create(&role).Error
}

func getRole(db database.Service, roleName string) uuid.UUID {
	role := new(models.Role)

	if err := db.UseGorm().Where("role = ?", role.Role).First(&role).Error; err != nil {
		panic("Unable to get role")
	}

	return role.Id
}
