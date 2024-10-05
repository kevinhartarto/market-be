package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email" gorm:"primaryKey"`
	Password  string    `json:"password"`
	Role      uuid.UUID `json:"role"`
	Verified  bool      `json:"verified"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	Id          uuid.UUID   `json:"id" gorm:"primaryKey"` // Referenced by Permission id
	Role        string      `json:"role"`
	Permissions Permissions `json:"permissions"`
}

type Permissions struct {
	Id       uuid.UUID `json:"id" gorm:"primaryKey"` // Referencing Role UUID
	IsViewer bool      `json:"viewer"`
	IsAdmin  bool      `json:"admin"`
	IsOwner  bool      `json:"owner"`
}
