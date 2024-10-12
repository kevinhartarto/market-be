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

	// Log in an identity.
	// returns an error if the identity cannot be identified.
	Login(c *fiber.Ctx) error

	// Deactivate an identity.
	// returns an error if the identity not found.
	Deactivate(c *fiber.Ctx) error

	// Create an admin identity.
	// returns an error if the identity cannot be created or already exists.
	CreateAdmin(c *fiber.Ctx) error

	// Create an user identity.
	// returns an error if the identity cannot be created or already exists.
	CreateUser(c *fiber.Ctx) error

	// Update an user
	// returns an error if the identity cannot be updated or does not exists.
	UpdateUser(c *fiber.Ctx) error

	// Get user by Id
	// returns an error if unable to find the user
	ShowUserDetails(c *fiber.Ctx) error

	// Create an identity role.
	// returns an error if similar role exits
	CreateRole(c *fiber.Ctx) error

	// Change user's role
	// return an error if it failed to update user's role
	ChangeRole(c *fiber.Ctx) error

	// Redefine role's permissions
	// return an error if the role not found
	ChangeRolePermissions(c *fiber.Ctx) error

	// Upgrade Role permissions to admin
	SetRoleToAdmin(c *fiber.Ctx) error

	DeleteRole(c *fiber.Ctx) error
}

var controllerInstance *identityController

type identityController struct {
	db    database.Service
	redis *redis.Client
}

// Temporary user for login credentials check
type tempUser struct {
	email    string
	password string
}

type roleChanger struct {
	user models.User
	role string
}

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

func (ic *identityController) Login(c *fiber.Ctx) error {
	tempUser := new(tempUser)
	user := new(models.User)

	if err := c.BodyParser(&tempUser); err != nil {
		return err
	}

	if err := ic.db.UseGorm().Where("email = ? and active", tempUser.email).First(&user).Error; err != nil {
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

	hashString := utils.HashString(user.Email)
	result := ic.redis.SetNX(c.Context(), hashString, tokenString, 0)
	fmt.Printf("Result from redis: %s", result)

	return c.JSON(fiber.Map{"token": tokenString})
}

func (ic *identityController) Deactivate(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	if err := ic.db.UseGorm().First(&role, role.Id).Error; err != nil {
		return err
	}

	return ic.db.UseGorm().Model(&role).Update("Active", false).Error
}

func (ic *identityController) CreateAdmin(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	user.Id = uuid.New()
	user.Role = getRole(db, "admin")
	user.Verified = false
	user.Active = true

	return ic.db.UseGorm().Create(&user).Error
}

func (ic *identityController) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	// Hash the password for the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Server Error")
	}

	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}

	user.Password = hashedPassword
	user.Role = getRole(db, "user")
	user.Verified = false
	user.Active = true

	return ic.db.UseGorm().Create(&user).Error
}

func (ic *identityController) UpdateUser(c *fiber.Ctx) error {
	user := new(models.Role)

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	return ic.db.UseGorm().Save(&user).Error
}

func (ic *identityController) ShowUserDetails(c *fiber.Ctx) error {
	userId := c.Query("id")
	user := new(models.User)

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request",
		})
	}

	return ic.db.UseGorm().First(&user, userId).Error
}

func (ic *identityController) CreateRole(c *fiber.Ctx) error {
	if err := c.BodyParser(&requestBody); err != nil {
		return err
	}

	role := requestBody.role
	permission := requestBody.permission

	// Handle UUID
	if role.Id == uuid.Nil {
		role.Id = uuid.New()
		permission.Id = role.Id
	} else {
		permission.Id = role.Id
	}

	if err := ic.db.UseGorm().Create(&role).Error; err != nil {
		return err
	}

	// Role created, now define Role's permissions
	if err := ic.db.UseGorm().Create(&permission).Error; err != nil {
		return err
	}

	// Role's permissions defined, now update it
	return ic.db.UseGorm().Model(&role).Update("permissions", permission).Error
}

func (ic *identityController) ChangeRole(c *fiber.Ctx) error {
	var tempRoleChanger roleChanger

	if err := c.BodyParser(&tempRoleChanger); err != nil {
		return err
	}

	role := new(models.Role)
	if err := ic.db.UseGorm().Where("name = ?", tempRoleChanger.role).Find(&role).Error; err != nil {
		return err
	}

	user := tempRoleChanger.user
	if err := ic.db.UseGorm().First(&user).Error; err != nil {
		return err
	}

	user.Role = role.Id

	return ic.db.UseGorm().Save(&user).Error
}

func (ic *identityController) ChangeRolePermissions(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	return ic.db.UseGorm().Save(&role).Error
}

func (ic *identityController) SetRoleToAdmin(c *fiber.Ctx) error {
	role := new(models.Role)
	roleId := c.Query("Id")

	if roleId == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	if err := ic.db.UseGorm().First(&role, roleId).Error; err != nil {
		return err
	}

	role.IsAdmin = true

	return ic.db.UseGorm().Save(&role).Error
}

func (ic *identityController) DeleteRole(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := ic.db.UseGorm().First(&role, role.Id).Error; err != nil {
		return err
	}

	return ic.db.UseGorm().Model(&role).Update("deprecated", true).Error
}

func (ic *identityController) GetRoleUUID(roleName string) uuid.UUID {
	role := new(models.Role)

	if err := ic.db.UseGorm().Where("role = ? ", roleName).First(&role).Error; err != nil {
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
