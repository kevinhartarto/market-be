package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

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
	UpdateAccount(c *fiber.Ctx) error

	// Get Account by Id
	// returns an error if unable to find the Account
	GetAccount(c *fiber.Ctx) error

	GetAllRoles(c *fiber.Ctx) error

	// load all roles from cache
	LoadRoles(ctx context.Context)

	// Create an identity role.
	// returns an error if similar role exits
	CreateRole(c *fiber.Ctx) error

	// Update a role
	// return an error if the role not found
	UpdateRole(c *fiber.Ctx) error
}

var (
	accountInstance *accountController

	account      models.Account
	roles        []models.Role
	role         models.Role
	loginCreds   loginCredentials
	affectedRows int64

	rolesKey   = "roles"
	unverified = "unverified"
	verified   = "verified"
)

type accountController struct {
	db    database.Service
	redis *redis.Client
}

type loginCredentials struct {
	email    string
	password string
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
	if err := c.BodyParser(&loginCreds); err != nil {
		return err
	}

	if err := ac.db.UseGorm().Where("email = ? and active", loginCreds.email).First(&account).Error; err != nil {
		return err
	}

	if err := utils.VerifyPassword(loginCreds.password, account.Password); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	tokenString := utils.GenerateJWT(account.Email, account.Role)
	if tokenString == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	hashString := utils.HashString(account.Email)
	result := ac.redis.SetNX(c.Context(), hashString, tokenString, 0)
	fmt.Printf("Result from redis: %s", result)

	return c.JSON(fiber.Map{"token": tokenString})
}

func (ac *accountController) CreateAccount(c *fiber.Ctx) error {
	if err := c.BodyParser(&account); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Hash the password for the Account
	hashedPassword, err := utils.HashPassword(account.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	account.Password = hashedPassword
	account.Role = getRole(c.Context(), *ac.redis, ac.db, unverified)

	if err := ac.db.UseGorm().Create(&account).Error; err != nil {
		return err
	}

	result, _ := json.Marshal(&account)
	return c.SendString(string(result))
}

func (ac *accountController) UpdateAccount(c *fiber.Ctx) error {
	var updateAccount struct {
		account     models.Account
		updateType  string
		updateValue bool
	}
	success := false

	if err := c.BodyParser(&updateAccount); err != nil {
		return err
	}

	switch updateAccount.updateType {
	case "update":
		affectedRows = ac.db.UseGorm().Save(&updateAccount.account).RowsAffected
	case "status":
		affectedRows = ac.db.UseGorm().Model(&updateAccount.account).Update("active", updateAccount.updateValue).RowsAffected
	case "verified":
		affectedRows = ac.db.UseGorm().Model(&updateAccount.account).Update("role", getRole(c.Context(), *ac.redis, ac.db, verified)).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&updateAccount.account)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (ac *accountController) GetAccount(c *fiber.Ctx) error {
	accountId := c.Query("id")
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
	if err := c.BodyParser(&role); err != nil {
		return err
	}

	// FIX THIS
	// if role.Id != uuid.Nil {
	// 	role.Id = uuid.New()
	// }

	if err := ac.db.UseGorm().Create(&role).Error; err != nil {
		return err
	}

	return c.SendString("Role " + role.Name + " created (" + role.Id.String() + ").")
}

func (ac *accountController) UpdateRole(c *fiber.Ctx) error {
	var updateRole struct {
		role        models.Role
		updateType  string
		updateValue bool
	}
	success := false

	if err := c.BodyParser(&updateRole); err != nil {
		return err
	}

	switch updateRole.updateType {
	case "update":
		affectedRows = ac.db.UseGorm().Save(&updateRole.role).RowsAffected
	case "admin":
		affectedRows = ac.db.UseGorm().Model(&updateRole.role).Update("is_admin", updateRole.updateValue).RowsAffected
	case "status":
		affectedRows = ac.db.UseGorm().Model(&updateRole.role).Update("deprecated", updateRole.updateValue).RowsAffected
	}

	// This is not a batch updates
	// Expect only 1 row changed
	if affectedRows == 1 {
		success = true
	}
	result, _ := json.Marshal(&updateRole.role)

	if success {
		return c.SendString(string(result))
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (ac *accountController) GetAllRoles(c *fiber.Ctx) error {
	if err := ac.redis.HGetAll(c.Context(), rolesKey).Scan(&roles); err != nil {
		fmt.Println("Roles not found in cache")
	}

	if reflect.ValueOf(&roles).IsZero() {
		ac.db.UseGorm().Where("deprecated is not true").Find(&roles)
	}

	result, _ := json.Marshal(roles)

	return c.SendString(string(result))
}

func (ac *accountController) LoadRoles(ctx context.Context) {
	ac.db.UseGorm().Where("deprecated is not true").Order("id asc").Find(&roles)
	ac.redis.Set(ctx, rolesKey, &roles, 0)
}

func getRole(ctx context.Context, redis redis.Client, db database.Service, roleName string) uuid.UUID {
	if err := redis.HGetAll(ctx, rolesKey).Scan(&roles); err != nil {
		fmt.Println("Roles not found in cache")
	}

	for _, tempRole := range roles {
		if tempRole.Name == roleName {
			return tempRole.Id
		}
	}

	if err := db.UseGorm().Where("role = ? and deprecated is not true", roleName).First(&role).Error; err != nil {
		panic("Unable to get role")
	}

	return role.Id
}
