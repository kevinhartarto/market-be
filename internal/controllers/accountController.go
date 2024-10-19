package controllers

import (
	"context"
	"encoding/json"
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

	// Create an Account identity.
	// returns an error if the identity cannot be created or already exists.
	CreateAccount(c *fiber.Ctx) error

	// Update an Account
	// returns an error if the identity cannot be updated or does not exists.
	UpdateAccount(c *fiber.Ctx, updateType string) error

	// Get Account by Id
	// returns an error if unable to find the Account
	GetAccount(c *fiber.Ctx) error

	// Get all roles
	GetAllRoles(ctx context.Context)

	// Create an identity role.
	// returns an error if similar role exits
	CreateRole(c *fiber.Ctx) error

	// Update a role
	// return an error if the role not found
	UpdateRole(c *fiber.Ctx, updateType string) error
}

var (
	accountInstance *accountController
	rolesKey        = "roles"
	unverified      = "unverified"
	verified        = "verified"
	admin           = "Admin"
)

type accountController struct {
	db    database.Service
	redis *redis.Client
}

// Temporary Account for login credentials check
type loginCredentials struct {
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
	loginCreds := new(loginCredentials)
	account := new(models.Account)

	if err := c.BodyParser(&loginCreds); err != nil {
		return err
	}

	if err := ac.db.UseGorm().Where("email = ? and active", loginCreds.email).First(&account).Error; err != nil {
		return err
	}

	if err := utils.VerifyPassword(loginCreds.password, account.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	tokenString := utils.GenerateJWT(account.Email, account.Role)
	if tokenString == "" {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	hashString := utils.HashString(account.Email)
	result := ac.redis.SetNX(c.Context(), hashString, tokenString, 0)
	fmt.Printf("Result from redis: %s", result)

	return c.JSON(fiber.Map{"token": tokenString})
}

func (ac *accountController) CreateAccount(c *fiber.Ctx) error {
	account := new(models.Account)

	if err := c.BodyParser(&account); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	}

	// Hash the password for the Account
	hashedPassword, err := utils.HashPassword(account.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Server Error")
	}

	if account.Id == uuid.Nil {
		account.Id = uuid.New()
	}

	account.Password = hashedPassword
	account.Role = getRole(c.Context(), *ac.redis, ac.db, unverified)
	account.Verified = false
	account.Active = true

	if err := ac.db.UseGorm().Create(&account).Error; err != nil {
		return err
	}

	result, _ := json.Marshal(&account)
	return c.SendString(string(result))
}

func (ac *accountController) UpdateAccount(c *fiber.Ctx, updateType string) error {
	account := new(models.Role)

	if err := c.BodyParser(&account); err != nil {
		return err
	}

	success := false
	dbConn := ac.db.UseGorm()
	var affectedRows int64

	switch updateType {
	case "update":
		affectedRows = dbConn.Save(&account).RowsAffected
	case "delete":
		affectedRows = dbConn.Model(&account).Update("Active", false).RowsAffected
	case "verified":
		var newRole = getRole(c.Context(), *ac.redis, ac.db, verified)
		affectedRows = dbConn.Model(&account).Update("role", newRole).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&account)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (ac *accountController) GetAccount(c *fiber.Ctx) error {
	accountId := c.Query("id")
	account := new(models.Account)

	if accountId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request",
		})
	}

	if err := ac.db.UseGorm().First(&account, accountId); err != nil {
		return c.SendString("error: Unable to find account")
	}

	result, _ := json.Marshal(&account)

	return c.SendString(string(result))
}

func (ac *accountController) CreateRole(c *fiber.Ctx) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	if role.Id != uuid.Nil {
		role.Id = uuid.New()
	}

	if err := ac.db.UseGorm().Create(&role).Error; err != nil {
		return err
	}

	return c.SendString("Role " + role.Name + " created (" + role.Id.String() + ").")
}

func (ac *accountController) UpdateRole(c *fiber.Ctx, updateType string) error {
	role := new(models.Role)

	if err := c.BodyParser(&role); err != nil {
		return err
	}

	success := false
	dbConn := ac.db.UseGorm()
	var affectedRows int64

	switch updateType {
	case "update":
		affectedRows = dbConn.Save(&role).RowsAffected
	case "upgrade":
		affectedRows = dbConn.Model(&role).Update("is_admin", true).RowsAffected
	case "delete":
		affectedRows = dbConn.Model(&role).Update("Active", false).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&role)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (ac *accountController) GetAllRoles(ctx context.Context) {
	var roles []models.Role

	ac.db.UseGorm().Where("is not deprecated").Order("id asc").Find(&roles)
	ac.redis.Set(ctx, rolesKey, &roles, 0)
}

func getRole(ctx context.Context, redis redis.Client, db database.Service, roleName string) uuid.UUID {
	var roles []models.Role
	role := new(models.Role)

	err := redis.HGetAll(ctx, rolesKey).Scan(&roles)
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
