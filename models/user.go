package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type UserDetail struct {
	UserID     uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"not null" json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `gorm:"not null" json:"last_name"`
	User       User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type UserRegitrationRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
