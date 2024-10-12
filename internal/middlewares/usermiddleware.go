package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kevinhartarto/market-be/internal/controllers"
	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/utils"
	"github.com/redis/go-redis/v9"
)

var SecretKey = []byte("e7185081-044a-4b23-ae05-95e18110607d")

type UserMiddleware struct {
	ctx   context.Context
	redis *redis.Client
}

func NewUserMiddleware(ctx context.Context, redisClient redis.Client) *UserMiddleware {
	return &UserMiddleware{
		ctx:   ctx,
		redis: &redisClient,
	}
}

type UserClaims struct {
	Email string    `json:"email"`
	Role  uuid.UUID `json:"role"`
	jwt.RegisteredClaims
}

func (um *UserMiddleware) Authorize(allowedRoles []string, db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token",
			})
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		key := utils.HashString(claims.Email)
		value, err := um.redis.Get(um.ctx, key).Result()
		if err != nil {
			return err
		}

		if value != tokenString {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		rolesUUID := controllers.GetRolesUUID(db)

		hasRole := false
		for _, roleUUID := range rolesUUID {
			if claims.Role == roleUUID {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: insufficient permissions",
			})
		}

		return c.Next()
	}
}
