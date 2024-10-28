package services

import (
	"log"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDatabaseOperationService_CreateUser(t *testing.T) {

	t.Run("successfully creates a user and user detail", func(t *testing.T) {
		user := &models.User{
			Email:    mocks.TestUserEmail,
			Password: mocks.TestUserPasswordHash,
		}
		userDetails := &models.UserDetail{
			FirstName: mocks.TestUserFirstName,
			LastName:  mocks.TestUserLastName,
		}

		err := DBOperationService.CreateUser(user, userDetails)
		require.NoError(t, err, "error should be nil when creating user and user details")

		// Verify the user was created
		var createdUser models.User
		err = DBOperationService.db.Where("email = ?", user.Email).First(&createdUser).Error
		require.NoError(t, err)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.Equal(t, user.Password, createdUser.Password)

		// Verify the user detail was created
		var createdUserDetail models.UserDetail
		err = DBOperationService.db.Where("user_id = ?", createdUser.ID).First(&createdUserDetail).Error
		require.NoError(t, err)
		assert.Equal(t, userDetails.FirstName, createdUserDetail.FirstName)
		assert.Equal(t, userDetails.LastName, createdUserDetail.LastName)

		sqlDB, err := DBOperationService.db.DB()
		if err != nil {
			log.Printf("Failed to connect to database for migrations: %v", err)
		}
		tests.DeleteTestData(sqlDB)
	})
}

func TestDatabaseOperationService_FindUserByEmail(t *testing.T) {
	t.Run("finds a user by email", func(t *testing.T) {
		user := &models.User{
			Email:    mocks.TestUserEmail,
			Password: mocks.TestUserPasswordHash,
		}
		userDetails := &models.UserDetail{
			FirstName: mocks.TestUserFirstName,
			LastName:  mocks.TestUserLastName,
		}

		err := DBOperationService.CreateUser(user, userDetails)
		require.NoError(t, err, "error should be nil when creating user and user details")

		// Find user by email
		foundUser, err := DBOperationService.FindUserByEmail(mocks.TestUserEmail)
		require.NoError(t, err)
		assert.Equal(t, user.Password, foundUser.Password)

		sqlDB, err := DBOperationService.db.DB()
		if err != nil {
			log.Printf("Failed to connect to database for migrations: %v", err)
		}
		tests.DeleteTestData(sqlDB)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		_, err := DBOperationService.FindUserByEmail("nonexistent@example.com")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestDatabaseOperationService_FindUserDetailsByUserID(t *testing.T) {

	t.Run("finds user details by user ID", func(t *testing.T) {
		user := &models.User{
			Email:    mocks.TestUserEmail,
			Password: mocks.TestUserPasswordHash,
		}
		userDetails := &models.UserDetail{
			UserID:    user.ID,
			FirstName: mocks.TestUserFirstName,
			LastName:  mocks.TestUserLastName,
		}

		err := DBOperationService.CreateUser(user, userDetails)
		require.NoError(t, err, "error should be nil when creating user and user details")

		// Find user details by user ID
		foundUserDetail, err := DBOperationService.FindUserDetailsByUserID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, userDetails.FirstName, foundUserDetail.FirstName)
		assert.Equal(t, userDetails.LastName, foundUserDetail.LastName)
	})

	t.Run("returns error when user details not found", func(t *testing.T) {
		_, err := DBOperationService.FindUserDetailsByUserID(9999) // non-existing user ID
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
