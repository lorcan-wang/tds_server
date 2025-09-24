package model

import "time"

type UserToken struct {
	ID           uint      `gorm:"primaryKey:autoIncrement"`
	UserID       string    `gorm:"type:uuid;not null;index"`
	AccessToken  string    `gorm:"type:text;not null"`
	RefreshToken string    `gorm:"type:text;not null"`
	ExpiresAt    time.Time `gorm:"not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
