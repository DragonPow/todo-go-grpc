package domain

import (
	"time"
)

type Task struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDone      bool      `json:"is_done"`
	DoneAt      time.Time `json:"done_at"`
	CreatedAt   time.Time `json:"created_at"`
	CreatorId   int32     `json:"creator_id"`
	TagsId      []int32   `json:"tags_id"`
}
