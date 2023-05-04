package model

import (
	"APM-server/pkg/auth"
	"time"

	"gorm.io/gorm"
)

// UserM 是数据库中 user 记录 struct 格式的映射.
type UserM struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Username  string         `gorm:"column:username"`
	Password  string         `gorm:"column:password;not null"`
	Email     string         `gorm:"column:email;not null;unique"`
	Status    int8           `gorm:"column:status;default:0"`
	AvatarUrl string         `gorm:"column:avatarUrl"`
	CreatedAt time.Time      `gorm:"column:createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deletedAt"`
}

// TableName 用来指定映射的 MySQL 表名.
func (u *UserM) TableName() string {
	return "user"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
