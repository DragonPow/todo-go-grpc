package internal

import (
	"context"
	"log"
	"testing"
	"time"
	"todo-go-grpc/app/response_handler"
	api "todo-go-grpc/app/task/api/task"
	"todo-go-grpc/app/task/domain"
	"todo-go-grpc/app/task/repository/mocks"
	user_service_api "todo-go-grpc/app/user/api"
	user_service_mock "todo-go-grpc/app/user/api/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	time_now         = time.Now()
	time_init        = time.Unix(0, 0).UTC()
	_taskRepoMock    = &mocks.TaskRepository{}
	_tagRepoMock     = &mocks.TagRepository{}
	_userServiceMock = &user_service_mock.UserHandlerClient{}
	_server          = server{
		taskRepo:    _taskRepoMock,
		tagRepo:     _tagRepoMock,
		userService: _userServiceMock,
	}
	ctx = context.Background()
)

func InitMock(t *testing.T) {
	// found user if id = 1
	_userServiceMock.On("Get", ctx, &user_service_api.GetReq{Id: 1}).Return(
		&user_service_api.User{
			Id:          1,
			Name:        "Thạch VN",
			Username:    "vungocthach",
			Password:    "123456",
			CreatedTime: timestamppb.New(time_now),
		}, nil)
	// not found user if id = 2
	_userServiceMock.On("Get", ctx, &user_service_api.GetReq{Id: 2}).Return(
		nil, response_handler.ResponseErrorNotFound(domain.ErrUserNotExists),
	)

	// fetch all wil return tag id : 1-> 4
	_tagRepoMock.On("FetchAll", ctx).Return(
		[]domain.Tag{
			{ID: 1, Value: "Value 1", Description: "Description 1", CreatedAt: time_now},
			{ID: 2, Value: "Value 2", Description: "Description 2", CreatedAt: time_now},
			{ID: 3, Value: "Value 3", Description: "Description 3", CreatedAt: time_now},
			{ID: 4, Value: "Value 4", Description: "Description 4", CreatedAt: time_now},
		}, nil,
	)
}

type userApiMock struct {
	user *user_service_api.User
	err  error
}

func TestList(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		tasks []domain.Task
		err   error
	}

	type testcase struct {
		expectIn  *api.ListReq
		expectOut *api.ListTask
		expectErr error
		data      *dataMock
	}

	taskApiMock := []*api.Task{
		{
			Id:          1,
			Name:        "Quét nhà",
			Description: "No description",
			IsDone:      false,
			Creator:     &api.User{Id: 1, Name: "Thạch VN", Username: "vungocthach"},
			Tags: []*api.Tag{
				{Id: 1, Value: "Value 1", Description: "Description 1"},
				{Id: 2, Value: "Value 2", Description: "Description 2"},
			},
			CreatedTime: timestamppb.New(time_now.Add(-time.Hour * 5)),
			DonedTime:   timestamppb.New(time_init),
		},
		{
			Id:          2,
			Name:        "Lau nhà",
			Description: "No description",
			IsDone:      true,
			Creator:     &api.User{Id: 1, Name: "Thạch VN", Username: "vungocthach"},
			Tags: []*api.Tag{
				{Id: 1, Value: "Value 1", Description: "Description 1"},
				{Id: 3, Value: "Value 3", Description: "Description 3"},
			},
			CreatedTime: timestamppb.New(time_now.Add(-time.Hour * 10)),
			DonedTime:   timestamppb.New(time_now),
		},
		{
			Id:          3,
			Name:        "Lau cửa",
			Description: "No description",
			IsDone:      false,
			Creator:     &api.User{Id: 1, Name: "Thạch VN", Username: "vungocthach"},
			Tags:        []*api.Tag{},
			CreatedTime: timestamppb.New(time_now),
			DonedTime:   timestamppb.New(time_init),
		},
	}
	taskModelMock := []domain.Task{

		{
			ID:          1,
			Name:        "Quét nhà",
			Description: "No description",
			IsDone:      false,
			CreatorId:   1,
			Tags: []domain.Tag{
				{ID: 1, Value: "Value 1", Description: "Description 1"},
				{ID: 2, Value: "Value 2", Description: "Description 2"},
			},
			CreatedAt: time_now.Add(-time.Hour * 5),
			DoneAt:    time_init,
		},
		{
			ID:          2,
			Name:        "Lau nhà",
			Description: "No description",
			IsDone:      true,
			CreatorId:   1,
			Tags: []domain.Tag{
				{ID: 1, Value: "Value 1", Description: "Description 1"},
				{ID: 3, Value: "Value 3", Description: "Description 3"},
			},
			CreatedAt: time_now.Add(-time.Hour * 10),
			DoneAt:    time_now,
		},
		{
			ID:          3,
			Name:        "Lau cửa",
			Description: "No description",
			IsDone:      false,
			CreatorId:   1,
			Tags:        []domain.Tag{},
			CreatedAt:   time_now,
			DoneAt:      time_init,
		},
	}

	testcases := []testcase{
		{
			// success: size 2 token 0, no name, no filter
			expectIn: &api.ListReq{
				PageSize:  2,
				PageToken: 0,
				Name:      "",
				Filter:    api.Filter_FILTER_UNSPECIFIED,
			},
			expectOut: &api.ListTask{
				Tasks: taskApiMock[:1],
			},
			data: &dataMock{
				tasks: taskModelMock[:1],
				err:   nil,
			},
		},
		{
			// success: size 2 token 1, no name, no filter
			expectIn: &api.ListReq{
				PageSize:  2,
				PageToken: 1,
				Name:      "",
				Filter:    api.Filter_FILTER_UNSPECIFIED,
			},
			expectOut: &api.ListTask{
				Tasks: taskApiMock[2:],
			},
			data: &dataMock{
				tasks: taskModelMock[2:],
				err:   nil,
			},
		},
		{
			// success: size 10 token 0, name: "Lau", no filter
			expectIn: &api.ListReq{
				PageSize:  10,
				PageToken: 0,
				Name:      "",
				Filter:    api.Filter_FILTER_UNSPECIFIED,
			},
			expectOut: &api.ListTask{
				Tasks: taskApiMock[1:2],
			},
			data: &dataMock{
				tasks: taskModelMock[1:2],
				err:   nil,
			},
		},
		{
			// success: size 10 token 0, no name, filter create desc
			expectIn: &api.ListReq{
				PageSize:  10,
				PageToken: 0,
				Name:      "",
				Filter:    api.Filter_TIME_CREATE_DESC,
			},
			expectOut: &api.ListTask{
				Tasks: []*api.Task{taskApiMock[2], taskApiMock[1], taskApiMock[0]},
			},
			data: &dataMock{
				tasks: []domain.Task{taskModelMock[2], taskModelMock[1], taskModelMock[0]},
				err:   nil,
			},
		},
	}

	for index, tc := range testcases {
		if tc.data == nil {
			_taskRepoMock.AssertNotCalled(t, "Fetch")
		} else {
			conditions := GetConditions(tc.expectIn)
			_taskRepoMock.On("Fetch", ctx, int32(1), tc.expectIn.PageToken, tc.expectIn.PageSize, conditions).
				Return(tc.data.tasks, tc.data.err)
		}

		actualOut, actualErr := _server.List(ctx, tc.expectIn)

		log.Printf("Output assert, testcase %v: %v", index+1, assert.Equal(t, tc.expectOut, actualOut))
		log.Printf("Error assert, testcase %v: %v", index+1, assert.ErrorIs(t, tc.expectErr, actualErr))
	}
	_taskRepoMock.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		task *domain.Task
		err  error
	}

	type testcase struct {
		expectIn  *api.GetReq
		expectOut *api.Task
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn: &api.GetReq{Id: 1},
			expectOut: &api.Task{
				Id:          1,
				Name:        "Quét nhà",
				Description: "No description",
				IsDone:      false,
				Creator:     &api.User{Id: 1, Name: "Thạch VN", Username: "vungocthach"},
				Tags: []*api.Tag{
					{Id: 1, Value: "Value 1", Description: "Description 1"},
					{Id: 2, Value: "Value 2", Description: "Description 2"},
				},
				CreatedTime: timestamppb.New(time_now),
				DonedTime:   timestamppb.New(time_init),
			},
			data: &dataMock{
				task: &domain.Task{
					ID:          1,
					Name:        "Quét nhà",
					Description: "No description",
					IsDone:      false,
					CreatedAt:   time_now,
					DoneAt:      time_init,
					CreatorId:   1,
					Tags: []domain.Tag{
						{ID: 1, Value: "Value 1", Description: "Description 1"},
						{ID: 2, Value: "Value 2", Description: "Description 2"},
					},
				},
				err: nil,
			},
		},
		{
			// fail: task not found
			expectIn:  &api.GetReq{Id: 2},
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrTaskNotExists),
			data: &dataMock{
				task: nil,
				err:  domain.ErrTaskNotExists,
			},
		},
		{
			// fail: user not found
			expectIn:  &api.GetReq{Id: 3},
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrUserNotExists),
			data: &dataMock{
				task: &domain.Task{
					ID:          1,
					Name:        "Vũ Ngọc Thạch",
					Description: "No description",
					IsDone:      false,
					CreatedAt:   time_now,
					DoneAt:      time_init,
					CreatorId:   2,
					Tags: []domain.Tag{
						{ID: 1, Value: "Value 1", Description: "Description 1"},
						{ID: 2, Value: "Value 2", Description: "Description 2"},
					},
				},
				err: nil,
			},
		},
	}

	for index, tc := range testcases {
		if tc.data == nil {
			_taskRepoMock.AssertNotCalled(t, "GetByID")
		} else {
			_taskRepoMock.On("GetByID", ctx, tc.expectIn.Id).Return(tc.data.task, tc.data.err)
		}

		actualOut, actualErr := _server.Get(ctx, tc.expectIn)

		log.Printf("Output assert, testcase %v: %v", index+1, assert.Equal(t, tc.expectOut, actualOut))
		log.Printf("Error assert, testcase %v: %v", index+1, assert.ErrorIs(t, tc.expectErr, actualErr))
	}
	_taskRepoMock.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	InitMock(t)
	// define testcase
	type dataMock struct {
		task *domain.Task
		err  error
	}

	type testcase struct {
		expectIn  *api.CreateReq
		expectOut *api.BasicTask
		expectErr error
		data      *dataMock
		userApi   *userApiMock
	}

	// create testcase
	testcases := []testcase{
		{
			// success
			expectIn: &api.CreateReq{
				Name:        "Dọn nhà",
				Description: "No comment",
				IsDone:      false,
				Tags:        []int32{1, 2},
			},
			expectOut: &api.BasicTask{
				Id:          1,
				Name:        "Dọn nhà",
				Description: "No comment",
				IsDone:      false,
				CreatorId:   1,
				TagsId:      []int32{1, 2},
				CreatedTime: timestamppb.New(time_now),
				DonedTime:   timestamppb.New(time_init),
			},
			expectErr: nil,
			data: &dataMock{
				task: &domain.Task{
					ID:          1,
					Name:        "Dọn nhà",
					Description: "No comment",
					IsDone:      false,
					CreatorId:   1,
					Tags: []domain.Tag{
						domain.Tag{ID: 1, Value: "1"},
						domain.Tag{ID: 2, Value: "2"},
					},
					DoneAt:    time_init,
					CreatedAt: time_now,
				},
			},
			userApi: &userApiMock{
				user: &user_service_api.User{
					Id:   1,
					Name: "Thạch",
				},
				err: nil,
			},
		},
		{
			// fail: tag not exists
			expectIn: &api.CreateReq{
				Name:        "Dọn nhà 2",
				Description: "No comment",
				IsDone:      false,
				Tags:        []int32{1, 3},
			},
			expectOut: nil,
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrTagNotExists),
			data: &dataMock{
				task: nil,
				err:  domain.ErrTagNotExists,
			},
			userApi: &userApiMock{
				user: &user_service_api.User{
					Id:   1,
					Name: "Thạch",
				},
				err: nil,
			},
		},
	}

	// test
	for _, tc := range testcases {
		if tc.userApi.err != nil {
			_taskRepoMock.AssertNotCalled(t, "Create")
		} else {
			if tc.data == nil {
				_taskRepoMock.AssertNotCalled(t, "Create")
			} else {
				// Tranfer data
				dataTransfer := &domain.Task{
					Name:        tc.expectIn.Name,
					Description: tc.expectIn.Description,
					IsDone:      tc.expectIn.IsDone,
					Tags:        []domain.Tag{},
				}
				for _, task_id := range tc.expectIn.Tags {
					dataTransfer.Tags = append(dataTransfer.Tags, domain.Tag{ID: task_id})
				}

				_taskRepoMock.On("Create", ctx, int32(1), dataTransfer).Return(tc.data.task, tc.data.err)
			}
		}

		actualOut, actualErr := _server.Create(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_taskRepoMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		task *domain.Task
		err  error
	}

	type testcase struct {
		expectIn  *api.UpdateReq
		expectOut *api.BasicTask
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn: &api.UpdateReq{
				Id:          1,
				NewTaskInfo: &api.BasicTask{Name: "Vũ Ngọc Thạch", Description: "No description", IsDone: false},
				TagsAdded:   []int32{1, 2},
				TagsDeleted: []int32{3, 4},
			},
			expectOut: &api.BasicTask{
				Id:          1,
				Name:        "Vũ Ngọc Thạch",
				Description: "No description",
				IsDone:      false,
				CreatorId:   1,
				TagsId:      []int32{1, 2, 5},
				CreatedTime: timestamppb.New(time_now),
				DonedTime:   timestamppb.New(time_init),
			},
			data: &dataMock{
				task: &domain.Task{
					ID:          1,
					Name:        "Vũ Ngọc Thạch",
					Description: "No description",
					IsDone:      false,
					CreatedAt:   time_now,
					DoneAt:      time_init,
					CreatorId:   1,
					Tags:        []domain.Tag{{ID: 1, Value: "1"}, {ID: 2, Value: "2"}, {ID: 5, Value: "5"}},
				},
				err: nil,
			},
		},
		{
			// fail: tag not found
			expectIn: &api.UpdateReq{
				Id:          2,
				NewTaskInfo: &api.BasicTask{Name: "Rửa chén", Description: "No description", IsDone: false},
				TagsAdded:   []int32{1, 6}, // fail for 6 not exists
				TagsDeleted: nil,
			},
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrTagNotExists),
			data:      nil,
		},
		{
			// fail: task name is exists
			expectIn: &api.UpdateReq{
				Id:          3,
				NewTaskInfo: &api.BasicTask{Name: "Lau nhà", Description: "No description", IsDone: false},
			},
			expectErr: response_handler.ResponseErrorAlreadyExists(domain.ErrTaskExists),
			data: &dataMock{
				task: nil,
				err:  domain.ErrTaskExists,
			},
		},
		{
			// fail: task not found
			expectIn: &api.UpdateReq{
				Id:          4,
				NewTaskInfo: &api.BasicTask{Name: "Quét nhà", Description: "No description", IsDone: true},
			},
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrTaskNotExists),
			data: &dataMock{
				task: nil,
				err:  domain.ErrTaskNotExists,
			},
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_taskRepoMock.AssertNotCalled(t, "Update")
		} else {
			data := transferBasicTaskToDomain(tc.expectIn.NewTaskInfo)
			_taskRepoMock.On("Update", ctx, tc.expectIn.Id, data, tc.expectIn.TagsAdded, tc.expectIn.TagsDeleted).
				Return(tc.data.task, tc.data.err)
		}

		actualOut, actualErr := _server.Update(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_taskRepoMock.AssertExpectations(t)
}

func TestDeleteMultiple(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		err error
	}

	type testcase struct {
		expectIn  *api.DeleteMultipleReq
		expectOut *emptypb.Empty
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn: &api.DeleteMultipleReq{
				TasksId: []int32{1, 2},
			},
			expectOut: &emptypb.Empty{},
			data: &dataMock{
				err: nil,
			},
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_taskRepoMock.AssertNotCalled(t, "Delete")
		} else {
			_taskRepoMock.On("Delete", ctx, tc.expectIn.TasksId).Return(tc.data.err)
		}

		actualOut, actualErr := _server.DeleteMultiple(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_taskRepoMock.AssertExpectations(t)
}
