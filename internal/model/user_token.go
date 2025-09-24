package model

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	ID           uint      `gorm:"primaryKey:autoIncrement"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index;unique"`
	AccessToken  string    `gorm:"type:text;not null"`
	RefreshToken string    `gorm:"type:text;not null"`
	ExpiresAt    time.Time `gorm:"not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
