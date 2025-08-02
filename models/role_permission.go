package models

type Role struct {
	ID       uint   `gorm:"primaryKey"`
	RoleName string `gorm:"unique;not null"`
}

type Permission struct {
	ID             uint   `gorm:"primaryKey"`
	PermissionName string `gorm:"unique;not null"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}
