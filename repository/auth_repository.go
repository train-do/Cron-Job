package repository

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"project/database"
	"project/domain"
	"time"
)

type AuthRepository struct {
	db        *gorm.DB
	cacher    database.Cacher
	secretKey string
}

func NewAuthRepository(db *gorm.DB, cacher database.Cacher, secretKey string) *AuthRepository {
	return &AuthRepository{db: db, cacher: cacher, secretKey: secretKey}
}

func (repo AuthRepository) Authenticate(user domain.User) (string, bool, error) {
	if err := repo.db.Where(user).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, errors.New("invalid username and/or password")
	}

	tokenData, signature := generateToken(user, repo.secretKey)
	if err := repo.cacher.Set(tokenData, signature); err != nil {
		return "", true, err
	}

	return fmt.Sprintf("%s.%s", tokenData, signature), true, nil

}

func generateToken(user domain.User, secretKey string) (string, string) {
	data := fmt.Sprintf("%d:%s:%d", user.ID, user.Role, time.Now().Unix())

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	tokenData := base64.URLEncoding.EncodeToString([]byte(data))
	return tokenData, signature
}
