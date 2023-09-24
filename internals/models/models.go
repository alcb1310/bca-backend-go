package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"primary key;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

type Company struct {
	Base
	Name      string `json:"name" gorm:"not null;index"`
	Employees uint   `json:"employees" gorm:"not null;default:1"`
	IsActive  bool   `json:"is_active" gorm:"not null;default:true"`
}
