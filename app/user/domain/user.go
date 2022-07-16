package domain

import "time"

type User struct {
	ID        int32     `form:"-" json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `form:"username" json:"username" binding:"required" gorm:"column:username;not null;unique"`
	Password  string    `form:"password" json:"password" binding:"required" gorm:"column:password;not null"`
	Name      string    `form:"name" json:"name" binding:"required" gorm:"column:name;not null"`
	CreatedAt time.Time `form:"-" json:"created_at" gorm:"column:created_at;"`
}
