package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      uuid.UUID `json:"role"`
	Verified  bool      `json:"verified"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name"`
	CanView     bool      `json:"can_view"`
	CanAdd      bool      `json:"can_add"`
	CanEdit     bool      `json:"can_edit"`
	CanDelete   bool      `json:"can_delete"`
	CanBuy      bool      `json:"can_buy"`
	CanWishlist bool      `json:"can_wishlist"`
	IsAdmin     bool      `json:"is_admin"`
	IsOwner     bool      `json:"is_owner"`
	Deprecated  bool      `json:"deprecated"`
}
