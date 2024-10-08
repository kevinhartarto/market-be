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
	Name        string      `json:"role"`
	Permissions Permissions `json:"permissions"`
	ExpiryTime  time.Time   `json:"expiry_time"`
	Deprecated  bool        `json:"deprecated"`
}

type Permissions struct {
	Id          uuid.UUID `json:"id" gorm:"primaryKey"` // Referencing Role UUID
	CanView     bool      `json:"can_view"`
	CanEdit     bool      `json:"can_edit"`
	CanBuy      bool      `json:"can_buy"`
	CanDelete   bool      `json:"can_delete"`
	CanWishlist bool      `json:"can_wishlist"`
	CanAdd      bool      `json:"can_add"`
	IsAdmin     bool      `json:"admin"`
	IsOwner     bool      `json:"owner"`
}
