package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type UserDetail struct {
	UserID     uint   `gorm:"column:user_id;primaryKey" json:"user_id"`
	FirstName  string `gorm:"column:firstname;not null;size:100" json:"first_name"`
	MiddleName string `gorm:"column:middlename;size:100" json:"middle_name"`
	LastName   string `gorm:"column:lastname;not null;size:100" json:"last_name"`
}

type UserRegitrationRequest struct {
	Email      string `json:"email" binding:"required,email"`
	FirstName  string `json:"first_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserCredentials struct {
	Email      string `json:"email" binding:"required,email"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Password   string `json:"password" binding:"required"`
}
