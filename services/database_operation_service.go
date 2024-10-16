package services

import (
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"gorm.io/gorm"
)

type IDatabaseOperationService interface {
	CreateUser(user *models.User, userDetail *models.UserDetail) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserDetailsByUserID(userID uint) (*models.UserDetail, error)
}

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

func (s *DatabaseOperationService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserDetailsByUserID finds the user details by their user ID
func (s *DatabaseOperationService) FindUserDetailsByUserID(userID uint) (*models.UserDetail, error) {
	var userDetails models.UserDetail
	if err := s.db.Where("user_id = ?", userID).First(&userDetails).Error; err != nil {
		return nil, err
	}
	return &userDetails, nil
}
