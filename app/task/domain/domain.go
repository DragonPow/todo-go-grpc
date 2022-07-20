package domain

import (
	"time"
)

type Task struct {
	ID          int32     `json:"id" form:"-" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" form:"name" gorm:"column:name;not null"`
	Description string    `json:"description" form:"description" gorm:"column:description"`
	IsDone      bool      `json:"is_done" form:"is_done" gorm:"column:is_done;default:false"`
	DoneAt      time.Time `json:"done_at" form:"-" gorm:"column:done_at"`
	CreatedAt   time.Time `json:"created_at" form:"-" gorm:"column:created_at"`
	CreatorId   int32     `form:"-"`
	Tags        []Tag     `json:"tags" form:"tags" gorm:"many2many:task_tags"`
}

type Tag struct {
	ID          int32     `form:"-" json:"id" gorm:"primaryKey;autoIncrement"`
	Value       string    `form:"value" json:"value" gorm:"column:value;not null;unique"`
	Description string    `form:"description" json:"description" gorm:"column:description"`
	CreatedAt   time.Time `form:"-" json:"created_at" gorm:"column:created_at"`
}
