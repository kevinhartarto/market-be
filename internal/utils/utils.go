package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("e7185081-044a-4b23-ae05-95e18110607d")

type UserClaims struct {
	Email string    `json:"email"`
	Role  uuid.UUID `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	// Generate a salt with a cost factor of 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(email string, role uuid.UUID) string {
	claims := UserClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}

	return tokenString
}

func HashRole(role string) string {

	objStr := fmt.Sprintf("%+v", role)
	data := []byte(objStr)
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		panic("Failed to hash role")
	}
	hash := hasher.Sum(nil)
	hashString := hex.EncodeToString(hash)

	return hashString
}
