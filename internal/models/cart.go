package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Content   []string  `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
