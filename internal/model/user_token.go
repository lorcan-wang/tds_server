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

// IsExpired reports whether the token expires within the given lead time.
func (t *UserToken) IsExpired(lead time.Duration) bool {
	if t == nil {
		return true
	}
	deadline := time.Now().Add(lead)
	return !deadline.Before(t.ExpiresAt)
}
