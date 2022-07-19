package domain

import (
	"time"
	tagDomain "todo-go-grpc/app/tag/domain"
	userDomain "todo-go-grpc/app/user/domain"
)

type Task struct {
	ID          int32     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"column:name;not null"`
	Description string    `json:"description" gorm:"column:description"`
	IsDone      bool      `json:"is_done" gorm:"column:is_done;default:false"`
	DoneAt      time.Time `json:"done_at" gorm:"column:done_at"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`

	CreatorId   int32
	UserCreator userDomain.User `json:"user_creator" gorm:"foreignKey:CreatorId;constrain:OnUpdate:NO ACTION,OnDelete:CASCADE"`

	Tags []tagDomain.Tag `json:"tags" gorm:"many2many:task_tags"`
}
