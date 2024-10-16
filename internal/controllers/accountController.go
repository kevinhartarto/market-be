package controllers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/models"
	"github.com/kevinhartarto/market-be/internal/utils"
	"github.com/redis/go-redis/v9"
)

// Identity Controller represent a service that handle all things related to Identity
type AccountController interface {

	// Log in an identity.
	// returns an error if the identity cannot be identified.
	Login(c *fiber.Ctx) error

	// Deactivate an identity.
	// returns an error if the identity not found.
	Deactivate(c *fiber.Ctx) error

	// Create an admin identity.
	// returns an error if the identity cannot be created or already exists.
	CreateAdmin(c *fiber.Ctx) error

	// Create an Account identity.
	// returns an error if the identity cannot be created or already exists.
	CreateAccount(c *fiber.Ctx) error

	// Update an Account
	// returns an error if the identity cannot be updated or does not exists.
	UpdateAccount(c *fiber.Ctx) error

	// Get Account by Id
	// returns an error if unable to find the Account
	ShowAccountDetails(c *fiber.Ctx) error

	// Create an identity role.
	// returns an error if similar role exits
	CreateRole(c *fiber.Ctx) error

	// Change Account's role
	// return an error if it failed to update Account's role
	ChangeRole(c *fiber.Ctx) error

	// Redefine role's permissions
	// return an error if the role not found
	ChangeRolePermissions(c *fiber.Ctx) error

	// Upgrade Role permissions to admin
	SetRoleToAdmin(c *fiber.Ctx) error

	DeleteRole(c *fiber.Ctx) error

	GetAllRoles(ctx context.Context)
}

var (
	accountInstance *accountController
	Roles           = "roles"
)

type accountController struct {
	db    database.Service
	redis *redis.Client
}

// Temporary Account for login credentials check
type tempAccount struct {
	email    string
	password string
}

type roleChanger struct {
	account models.Account
	role    string
}

func NewAccountController(db database.Service, redis *redis.Client) *accountController {

	if accountInstance != nil {
		return accountInstance
	}

	accountInstance = &accountController{
		db:    db,
		redis: redis,
	}

	return accountInstance
}

func (ac *accountController) Login(c *fiber.Ctx) error {
	tempAccount := new(tempAccount)
	Account := new(models.Account)

	if err := c.BodyParser(&tempAccount); err != nil {
		return err
	}

	if err := ac.db.UseGorm().Where("email = ? and active", tempAccount.email).First(&Account).Error; err != nil {
		return err
	}

	if err := utils.VerifyPassword(tempAccount.password, Account.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	tokenString := utils.GenerateJWT(Account.Email, Account.Role)
	if tokenString == "" {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	hashString := utils.HashString(Account.Email)
	result := ac.redis.SetNX(c.Context(), hashString, tokenString, 0)
	fmt.Printf("Result from redis: %s", result)

	return c.JSON(fiber.Map{"token": tokenString})
}

func (ac *accountController) Deactivate(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	if err := ac.db.UseGorm().First(&role, role.Id).Error; err != nil {
		return err
	}

	return ac.db.UseGorm().Model(&role).Update("Active", false).Error
}

func (ac *accountController) CreateAdmin(c *fiber.Ctx) error {
	Account := new(models.Account)

	if err := c.BodyParser(Account); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	Account.Id = uuid.New()
	Account.Role = getRole(c.Context(), *ac.redis, ac.db, "admin")
	Account.Verified = false
	Account.Active = true

	return ac.db.UseGorm().Create(&Account).Error
}

func (ac *accountController) CreateAccount(c *fiber.Ctx) error {
	Account := new(models.Account)

	if err := c.BodyParser(&Account); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	// Hash the password for the Account
	hashedPassword, err := utils.HashPassword(Account.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Server Error")
	}

	if Account.Id == uuid.Nil {
		Account.Id = uuid.New()
	}

	Account.Password = hashedPassword
	Account.Role = getRole(c.Context(), *ac.redis, ac.db, "normal")
	Account.Verified = false
	Account.Active = true

	return ac.db.UseGorm().Create(&Account).Error
}

func (ac *accountController) UpdateAccount(c *fiber.Ctx) error {
	Account := new(models.Role)

	if err := c.BodyParser(&Account); err != nil {
		return err
	}

	return ac.db.UseGorm().Save(&Account).Error
}

func (ac *accountController) ShowAccountDetails(c *fiber.Ctx) error {
	AccountId := c.Query("id")
	Account := new(models.Account)

	if AccountId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request",
		})
	}

	return ac.db.UseGorm().First(&Account, AccountId).Error
}

func (ac *accountController) CreateRole(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	if role.Id != uuid.Nil {
		role.Id = uuid.New()
	}

	return ac.db.UseGorm().Create(&role).Error
}

func (ac *accountController) ChangeRole(c *fiber.Ctx) error {
	var tempRoleChanger roleChanger

	if err := c.BodyParser(&tempRoleChanger); err != nil {
		return err
	}

	role := new(models.Role)
	if err := ac.db.UseGorm().Where("name = ?", tempRoleChanger.role).Find(&role).Error; err != nil {
		return err
	}

	Account := tempRoleChanger.account
	if err := ac.db.UseGorm().First(&Account).Error; err != nil {
		return err
	}

	Account.Role = role.Id

	return ac.db.UseGorm().Save(&Account).Error
}

func (ac *accountController) ChangeRolePermissions(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	return ac.db.UseGorm().Save(&role).Error
}

func (ac *accountController) SetRoleToAdmin(c *fiber.Ctx) error {
	role := new(models.Role)
	roleId := c.Query("Id")

	if roleId == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	if err := ac.db.UseGorm().First(&role, roleId).Error; err != nil {
		return err
	}

	role.IsAdmin = true

	return ac.db.UseGorm().Save(&role).Error
}

func (ac *accountController) DeleteRole(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := ac.db.UseGorm().First(&role, role.Id).Error; err != nil {
		return err
	}

	return ac.db.UseGorm().Model(&role).Update("deprecated", true).Error
}

func (ac *accountController) GetAllRoles(ctx context.Context) {
	var roles []models.Role

	ac.db.UseGorm().Where("is not deprecated").Order("id asc").Find(&roles)
	ac.redis.Set(ctx, Roles, &roles, 0)
}

func getRole(ctx context.Context, redis redis.Client, db database.Service, roleName string) uuid.UUID {
	var roles []models.Role
	role := new(models.Role)

	err := redis.HGetAll(ctx, Roles).Scan(&roles)
	if err != nil {
		fmt.Println("Roles not found in cache")
	}

	for _, tempRole := range roles {
		if tempRole.Name == roleName {
			return tempRole.Id
		}
	}

	if err := db.UseGorm().Where("role = ? and is not deprecated", roleName).First(&role).Error; err != nil {
		panic("Unable to get role")
	}

	return role.Id
}
