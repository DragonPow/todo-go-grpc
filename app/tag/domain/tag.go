package domain

import "time"

type Tag struct {
	ID          int32     `form:"-" json:"id" gorm:"primaryKey;autoIncrement"`
	Value       string    `form:"value" json:"value" gorm:"column:value;not null;unique"`
	Description string    `form:"description" json:"description" gorm:"column:description"`
	CreatedAt   time.Time `form:"-" json:"created_at" gorm:"column:created_at"`
}
