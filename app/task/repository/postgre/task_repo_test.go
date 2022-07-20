package postgre

import (
	"context"
	"reflect"
	"testing"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/task/domain"
	"todo-go-grpc/app/task/repository"

	"gorm.io/gorm"
)

func TestNewTaskRepository(t *testing.T) {
	type args struct {
		conn dbservice.Database
	}
	tests := []struct {
		name string
		args args
		want repository.TaskRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskRepository(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchUserByIds(t *testing.T) {
	type args struct {
		ctx context.Context
		ids []int32
		db  *gorm.DB
	}
	tests := []struct {
		name      string
		args      args
		wantTasks []domain.Task
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTasks, err := SearchUserByIds(tt.args.ctx, tt.args.ids, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchUserByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTasks, tt.wantTasks) {
				t.Errorf("SearchUserByIds() = %v, want %v", gotTasks, tt.wantTasks)
			}
		})
	}
}

func Test_taskRepository_Fetch(t *testing.T) {
	type args struct {
		ctx        context.Context
		user_id    int32
		offset     int32
		number     int32
		conditions map[string]any
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		want    []domain.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Fetch(tt.args.ctx, tt.args.user_id, tt.args.offset, tt.args.number, tt.args.conditions)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskRepository.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRepository_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		want    *domain.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRepository_Create(t *testing.T) {
	type args struct {
		ctx        context.Context
		creator_id int32
		info       *domain.Task
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		want    *domain.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Create(tt.args.ctx, tt.args.creator_id, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRepository_Update(t *testing.T) {
	type args struct {
		ctx         context.Context
		id          int32
		new_info    *domain.Task
		tags_add    []int32
		tags_remove []int32
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		want    *domain.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Update(tt.args.ctx, tt.args.id, tt.args.new_info, tt.args.tags_add, tt.args.tags_remove)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskRepository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRepository_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		ids []int32
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Delete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskRepository_IsExists(t *testing.T) {
	type args struct {
		ctx context.Context
		ids int32
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.IsExists(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.IsExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("taskRepository.IsExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRepository_GetByUserId(t *testing.T) {
	type args struct {
		ctx     context.Context
		user_id int32
	}
	tests := []struct {
		name    string
		tr      *taskRepository
		args    args
		want    []int32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.GetByUserId(tt.args.ctx, tt.args.user_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepository.GetByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskRepository.GetByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}
