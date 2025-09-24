package repository

import (
	"tds_server/internal/data"
	"tds_server/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo() *TokenRepo {
	return &TokenRepo{db: data.DB}
}

// 保存（创建或更新）用户token
func (repo *TokenRepo) Save(userID uuid.UUID, accessToken string, refreshToken string, expiresIn time.Duration) error {
	expiresAt := time.Now().Add(expiresIn * time.Second)
	token := model.UserToken{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	return repo.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token", "expires_at"}),
		},
	).Create(&token).Error
}
