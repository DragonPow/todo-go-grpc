package postgre

import (
	"context"
	"errors"
	"todo-go-grpc/app/dbservice"
	tagDomain "todo-go-grpc/app/tag/domain"
	"todo-go-grpc/app/task/domain"
	"todo-go-grpc/app/task/repository"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type taskRepository struct {
	Conn dbservice.Database
}

func NewTaskRepository(conn dbservice.Database) repository.TaskRepository {
	return &taskRepository{
		Conn: conn,
	}
}

func SearchUserByIds(ctx context.Context, ids []int32, db *gorm.DB) (tasks []domain.Task, err error) {
	if err = db.Where("id IN ?", ids).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *taskRepository) Fetch(ctx context.Context, user_id int32, offset int32, number int32, conditions map[string]any) ([]domain.Task, error) {
	var tasks []domain.Task
	var queryString string
	tx := t.Conn.Db.Preload("UserCreator").Preload("Tags")
	queryArgs := []interface{}{}

	// Check condition and add to queryString
	if value, ok := conditions["name"]; ok {
		queryString += "name LIKE ?"
		queryArgs = append(queryArgs, "%"+value.(string)+"%")
	}
	// if tags, ok := conditions["tags"]; ok && tags != nil {
	// 	if queryString != "" {
	// 		queryString += " AND "
	// 	}

	// 	queryString += "tag_id IN ?"
	// 	queryArgs = append(queryArgs, tags.([]int32))
	// }

	if queryString != "" {
		tx = tx.Where(queryString, queryArgs...)
	}

	// Set order
	if filter, ok := conditions["filter"]; ok && filter != nil {
		switch filter {
		case "TIME_CREATE_ASC":
			tx = tx.Order("created_at asc")
		case "TIME_CREATE_DESC":
			tx = tx.Order("created_at desc")
		}
	} else {
		tx = tx.Order("id asc")
	}

	if err := tx.Limit(int(number)).Offset(int(offset)).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskRepository) GetByID(ctx context.Context, id int32) (*domain.Task, error) {
	var task domain.Task
	if err := t.Conn.Db.Preload("UserCreator").Preload("Tags").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTaskNotExists
		}

		return nil, err
	}

	return &task, nil
}

func (t *taskRepository) Create(ctx context.Context, creator_id int32, info *domain.Task) (*domain.Task, error) {
	info.CreatorId = creator_id
	if err := t.Conn.Db.Create(&info).Error; err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && errors.Is(err, pgError) {
			if pgError.Code == "23503" {
				return nil, domain.ErrTagNotExists
			}
		}
		return nil, err
	}

	// Add tags to task
	if err := t.Conn.Db.Model(&info).Association("Tags").Append(&info.Tags); err != nil {
		return nil, err
	}

	// if err := tx.Model(&new_task).Association("UserCreator").Replace(&new_task.UserCreator); err != nil {
	// 	return domain.Task{}, err
	// }

	return info, nil
}

func (t *taskRepository) Update(ctx context.Context, id int32, new_info *domain.Task, tags_add []int32, tags_remove []int32) (*domain.Task, error) {
	new_task_map := map[string]any{}
	new_task_map["name"] = new_info.Name
	new_task_map["description"] = new_info.Description
	new_task_map["is_done"] = new_info.IsDone
	new_task_map["creator_id"] = new_info.CreatorId

	// Update information
	if err := t.Conn.Db.First(&new_info, id).Updates(new_task_map).Error; err != nil {
		return nil, err
	}

	// Update tags
	tranferIdToTag := func(ids []int32) (tags []tagDomain.Tag) {
		for _, id := range ids {
			tags = append(tags, tagDomain.Tag{ID: id})
		}
		return tags
	}

	if err := t.Conn.Db.Model(&new_info).Association("Tags").Append(tranferIdToTag(tags_add)); err != nil {
		return nil, err
	}
	if err := t.Conn.Db.Model(&new_info).Association("Tags").Delete(tranferIdToTag(tags_remove)); err != nil {
		return nil, err
	}

	// TODO: implement new_info here
	return new_info, nil
}

func (t *taskRepository) Delete(ctx context.Context, ids []int32) error {
	if len(ids) == 0 {
		return nil
	}

	// Delete tasks
	// Add Select("Tags") to delete association of task and tag
	if err := t.Conn.Db.Select("Tags").Delete(&domain.Task{}, ids).Error; err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) IsExists(ctx context.Context, ids int32) (bool, error) {
	// Find by IDs
	var tasks []domain.Task
	if err := t.Conn.Db.Where("Id IN ?", ids).Find(&tasks).Error; err != nil {
		return false, err
	}

	// Check length of tasks found is equal ids
	if len(tasks) != 1 {
		return false, nil
	}

	return true, nil
}
