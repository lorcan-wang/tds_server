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

// Save creates or updates a user token. Save 方法会创建或更新用户 token。
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

// GetByUserID retrieves a user token by ID. GetByUserID 根据用户 ID 查询 token。
func (repo *TokenRepo) GetByUserID(userID uuid.UUID) (*model.UserToken, error) {
	var token model.UserToken
	if err := repo.db.Where("user_id = ?", userID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
