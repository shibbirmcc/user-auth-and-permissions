package services

import (
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"gorm.io/gorm"
)

type DatabaseOperationService struct {
	db *gorm.DB
}

func NewDatabaseOperationService(db *gorm.DB) *DatabaseOperationService {
	return &DatabaseOperationService{db: db}
}

func (s *DatabaseOperationService) CreateUser(user *models.User, userDetail *models.UserDetail) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userDetail.UserID = user.ID
		return tx.Create(userDetail).Error
	})
}
