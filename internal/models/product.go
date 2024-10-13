package models

import (
	"time"

	"github.com/google/uuid"
)

type Brand struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	OnSale    bool      `json:"on_sale"`
	Active    bool      `json:"active"`
	Owner     uuid.UUID `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdateBy  uuid.UUID `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Featured    bool      `json:"featured"`
	Active      bool      `json:"active"`
	Owner       uuid.UUID `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateBy    uuid.UUID `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name"`
	Image       []string  `json:"images"`
	Price       int       `json:"price"`
	Colour      []string  `json:"colours"`
	Brand       uuid.UUID `json:"brand"`
	Categories  []string  `json:"categories"`
	Size        []string  `json:"size"`
	OnSale      bool      `json:"on_sale"`
	SalePrice   int       `json:"sale_price"`
	SalePercent int       `json:"sale_percent"`
	Stock       int       `json:"stock"`
	IsNew       bool      `json:"is_new"`
	Description string    `json:"description"`
	Owner       uuid.UUID `json:"owner"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateBy    uuid.UUID `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}
