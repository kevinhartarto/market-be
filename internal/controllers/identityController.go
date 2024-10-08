package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"github.com/kevinhartarto/market-be/internal/utils"
	"github.com/redis/go-redis/v9"
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
	db    database.Service
	redis *redis.Client
}

type tempUser struct {
	email    string
	password string
}

var (
	controllerInstance *identityController
)
var SecretKey = []byte("SecretKey")

func NewIdentityController(db database.Service, redis *redis.Client) IdentityController {

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

	if err := db.UseGorm().Where("email = ? and active", tempUser.email).First(&user).Error; err != nil {
		return err
	}

	err := utils.VerifyPassword(tempUser.password, user.Password).Error
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	tokenString := utils.GenerateJWT(user.Email, user.Role)
	if tokenString == "" {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	hashString := utils.HashRole(user.Role.String())
	result := ic.redis.SetNX(c.Context(), hashString, tokenString, 0)
	fmt.Printf("Result from redis: %s", result)

	return c.JSON(fiber.Map{"token": tokenString})
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

func (ic *identityController) GetRoleUUID(roleName string, db database.Service) uuid.UUID {
	role := new(models.Role)

	if err := db.UseGorm().Where("role = ? ", roleName).First(&role).Error; err != nil {
		panic("Unable to get role")
	}

	return role.Id
}

func GetRolesUUID(db database.Service) uuid.UUIDs {
	var (
		roles []models.Role
		UUIDs uuid.UUIDs
	)

	if err := db.UseGorm().Find(&roles).Error; err != nil {
		panic("Unable to get roles")
	}

	for _, role := range roles {
		UUIDs = append(UUIDs, role.Id)
	}

	return UUIDs
}

func getRole(db database.Service, roleName string) uuid.UUID {
	role := new(models.Role)

	if err := db.UseGorm().Where("role = ? and is not deprecated", roleName).First(&role).Error; err != nil {
		panic("Unable to get role")
	}

	return role.Id
}
