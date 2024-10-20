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
	Verified  bool      `json:"verified" gorm:"default:false"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name"`
	CanView     bool      `json:"can_view" gorm:"default:true"`
	CanAdd      bool      `json:"can_add" gorm:"default:false"`
	CanEdit     bool      `json:"can_edit" gorm:"default:false"`
	CanDelete   bool      `json:"can_delete" gorm:"default:false"`
	CanBuy      bool      `json:"can_buy" gorm:"default:false"`
	CanWishlist bool      `json:"can_wishlist" gorm:"default:false"`
	IsAdmin     bool      `json:"is_admin" gorm:"default:false"`
	IsOwner     bool      `json:"is_owner" gorm:"default:false"`
	Deprecated  bool      `json:"deprecated" gorm:"default:false"`
}
