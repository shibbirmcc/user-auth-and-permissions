package utils

import (
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/stretchr/testify/assert"
)

func TestMarshalObject_Success_User(t *testing.T) {
	// Arrange: Create a User object
	user := models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "securePassword123",
	}

	// Act: Marshal the User object
	jsonBytes, err := MarshalObject(user)

	// Assert
	assert.NoError(t, err)
	assert.JSONEq(t, `{"ID":1,"email":"test@example.com","password":"securePassword123"}`, string(jsonBytes))
}

func TestMarshalObject_Success_UserDetail(t *testing.T) {
	// Arrange: Create a UserDetail object
	userDetail := models.UserDetail{
		UserID:     1,
		FirstName:  "John",
		MiddleName: "M",
		LastName:   "Doe",
	}

	// Act: Marshal the UserDetail object
	jsonBytes, err := MarshalObject(userDetail)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"user_id":1,"first_name":"John","middle_name":"M","last_name":"Doe"}`, string(jsonBytes))
}

func TestMarshalObject_Success_UserRegistrationRequest(t *testing.T) {
	// Arrange: Create a UserRegistrationRequest object
	registrationRequest := models.UserRegitrationRequest{
		Email:      "test@example.com",
		FirstName:  "John",
		MiddleName: "M",
		LastName:   "Doe",
	}

	// Act: Marshal the UserRegistrationRequest object
	jsonBytes, err := MarshalObject(registrationRequest)

	// Assert
	assert.NoError(t, err)
	assert.JSONEq(t, `{"email":"test@example.com","first_name":"John","middle_name":"M","last_name":"Doe"}`, string(jsonBytes))
}

func TestMarshalObject_Success_LoginRequest(t *testing.T) {
	// Arrange: Create a LoginRequest object
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "securePassword123",
	}

	// Act: Marshal the LoginRequest object
	jsonBytes, err := MarshalObject(loginRequest)

	// Assert
	assert.NoError(t, err)
	assert.JSONEq(t, `{"email":"test@example.com","password":"securePassword123"}`, string(jsonBytes))
}

func TestMarshalObject_Success_UserCredentials(t *testing.T) {
	// Arrange: Create a UserCredentials object
	credentials := models.UserCredentials{
		Email:      "test@example.com",
		FirstName:  "John",
		MiddleName: "M",
		LastName:   "Doe",
		Password:   "securePassword123",
	}

	// Act: Marshal the UserCredentials object
	jsonBytes, err := MarshalObject(credentials)

	// Assert
	assert.NoError(t, err)
	assert.JSONEq(t, `{"email":"test@example.com","first_name":"John","middle_name":"M","last_name":"Doe","password":"securePassword123"}`, string(jsonBytes))
}

func TestMarshalObject_NilInput(t *testing.T) {
	// Act: Attempt to marshal a nil input
	jsonBytes, err := MarshalObject(nil)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "input cannot be nil")
	assert.Nil(t, jsonBytes)
}

func TestMarshalObject_InvalidInput(t *testing.T) {
	// Arrange: Use an unsupported type for marshaling
	invalidInput := func() {}

	// Act: Attempt to marshal the invalid input
	jsonBytes, err := MarshalObject(invalidInput)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
	assert.Nil(t, jsonBytes)
}

func TestMarshalObject_EmptyStruct(t *testing.T) {
	emptyStruct := struct{}{}

	jsonBytes, err := MarshalObject(emptyStruct)

	assert.NoError(t, err)
	assert.Equal(t, `{}`, string(jsonBytes))
}

func TestMarshalObject_NestedStruct(t *testing.T) {
	nestedStruct := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Meta struct {
			Info string `json:"info"`
		} `json:"meta"`
	}{
		ID:   1,
		Name: "Test",
		Meta: struct {
			Info string `json:"info"`
		}{
			Info: "Nested data",
		},
	}

	jsonBytes, err := MarshalObject(nestedStruct)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"id":1,"name":"Test","meta":{"info":"Nested data"}}`, string(jsonBytes))
}

func TestMarshalObject_StructWithPointers(t *testing.T) {
	name := "John"
	pointerStruct := struct {
		ID   *int    `json:"id"`
		Name *string `json:"name"`
	}{
		ID:   nil,
		Name: &name,
	}

	jsonBytes, err := MarshalObject(pointerStruct)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"id":null,"name":"John"}`, string(jsonBytes))
}

func TestMarshalObject_LargeData(t *testing.T) {
	largeData := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		largeData = append(largeData, i)
	}

	jsonBytes, err := MarshalObject(largeData)

	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), "999")
}

func TestMarshalObject_UnsupportedType(t *testing.T) {
	unsupported := make(chan int)

	jsonBytes, err := MarshalObject(unsupported)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
	assert.Nil(t, jsonBytes)
}
