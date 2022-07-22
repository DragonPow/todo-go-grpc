package internal

import (
	"context"
	"testing"
	"time"
	"todo-go-grpc/app/response_handler"
	api "todo-go-grpc/app/task/api/task"
	"todo-go-grpc/app/task/domain"
	"todo-go-grpc/app/task/repository/mocks"
	user_service_api "todo-go-grpc/app/user/api"
	user_service_mock "todo-go-grpc/app/user/api/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T) {
	_taskRepoMock := &mocks.TaskRepository{}
	_tagRepoMock := &mocks.TagRepository{}
	_userServiceMock := &user_service_mock.UserHandlerClient{}
	_server := server{
		taskRepo:    _taskRepoMock,
		tagRepo:     _tagRepoMock,
		userService: _userServiceMock,
	}
	ctx := context.Background()

	type dataMock struct {
		task *domain.Task
		err  error
	}

	type userApiMock struct {
		user *user_service_api.User
		err  error
	}

	type testcase struct {
		expectIn  *api.CreateReq
		expectOut *api.BasicTask
		expectErr error
		data      *dataMock
		userApi   *userApiMock
	}

	time_now := time.Now()
	time_init := time.Unix(0, 0).UTC()
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
		// {
		// 	// fail: user not exists
		// 	expectIn: &api.CreateReq{
		// 		Name:        "Dọn nhà 3",
		// 		Description: "No comment",
		// 		IsDone:      false,
		// 		Tags:        []int32{1, 3},
		// 	},
		// 	expectOut: nil,
		// 	expectErr: response_handler.ResponseErrorNotFound(domain.ErrUserNotExists),
		// 	data: &dataMock{
		// 		task: nil,
		// 		err:  domain.ErrTagNotExists,
		// 	},
		// 	userApi: &userApiMock{
		// 		user: nil,
		// 		err:  domain.ErrUserNotExists,
		// 	},
		// },
	}

	for _, tc := range testcases {
		_userServiceMock.On("Get", ctx, &user_service_api.GetReq{Id: 1}).Return(tc.userApi.user, tc.userApi.err)
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

// func TestUpdate(t *testing.T) {
// 	_repoMock := &mocks.UserRepository{}
// 	_server := server{
// 		repo: _repoMock,
// 	}
// 	ctx := context.Background()

// 	type dataMock struct {
// 		user *domain.User
// 		err  error
// 	}

// 	type testcase struct {
// 		expectIn  *api.UpdateReq
// 		expectOut *api.User
// 		expectErr error
// 		data      *dataMock
// 	}

// 	created_time := time.Now()
// 	testcases := []testcase{
// 		{
// 			// success
// 			expectIn: &api.UpdateReq{
// 				Id: 1,
// 				NewUserInfor: &api.User{
// 					Name:        "Vũ Ngọc Thạch",
// 					Username:    "vungocthach",
// 					Password:    "123456",
// 					CreatedTime: timestamppb.New(created_time),
// 				},
// 			},
// 			expectOut: &api.User{
// 				Id:          1,
// 				Name:        "Vũ Ngọc Thạch",
// 				Username:    "vungocthach",
// 				Password:    "123456",
// 				CreatedTime: timestamppb.New(created_time),
// 			},
// 			expectErr: nil,
// 			data: &dataMock{
// 				user: &domain.User{
// 					ID:        1,
// 					Name:      "Vũ Ngọc Thạch",
// 					Username:  "vungocthach",
// 					Password:  "123456",
// 					CreatedAt: created_time,
// 				},
// 				err: nil,
// 			},
// 		},
// 		{
// 			// fail: duplicate username
// 			expectIn: &api.UpdateReq{
// 				Id: 1,
// 				NewUserInfor: &api.User{
// 					Name:     "Minh Nhực",
// 					Username: "minhnhuc",
// 					Password: "123456",
// 				},
// 			},
// 			expectOut: nil,
// 			expectErr: response_service.ResponseErrorAlreadyExists(domain.ErrUserNameIsExists),
// 			data: &dataMock{
// 				user: nil,
// 				err:  domain.ErrUserNameIsExists,
// 			},
// 		},
// 		{
// 			// fail: id not exists
// 			expectIn: &api.UpdateReq{
// 				Id:           2,
// 				NewUserInfor: &api.User{},
// 			},
// 			expectOut: nil,
// 			expectErr: response_service.ResponseErrorNotFound(domain.ErrUserNotExists),
// 			data:      nil,
// 		},
// 	}

// 	// If id = 1, found user, else if = 2, not found
// 	_repoMock.On("GetByID", ctx, mock.MatchedBy(func(id int32) bool { return id == 1 })).Return(nil, nil)
// 	_repoMock.On("GetByID", ctx, mock.MatchedBy(func(id int32) bool { return id == 2 })).Return(nil, domain.ErrUserNotExists)

// 	for _, tc := range testcases {
// 		if tc.data == nil {
// 			_repoMock.AssertNotCalled(t, "Update")
// 		} else {
// 			_repoMock.On("Update", ctx, tc.expectIn.Id,
// 				&domain.User{
// 					// ID:        tc.expectIn.Id,
// 					Name:      tc.expectIn.NewUserInfor.Name,
// 					Username:  tc.expectIn.NewUserInfor.Username,
// 					Password:  tc.expectIn.NewUserInfor.Password,
// 					CreatedAt: tc.expectIn.NewUserInfor.CreatedTime.AsTime(),
// 				},
// 			).Return(tc.data.user, tc.data.err)
// 		}

// 		actualOut, actualErr := _server.Update(ctx, tc.expectIn)

// 		assert.Equal(t, tc.expectOut, actualOut)
// 		assert.ErrorIs(t, tc.expectErr, actualErr)
// 	}
// 	_repoMock.AssertExpectations(t)
// }

// func TestDelete(t *testing.T) {
// 	_repoMock := &mocks.UserRepository{}
// 	_server := server{
// 		repo: _repoMock,
// 	}
// 	ctx := context.Background()

// 	type dataMock struct {
// 		err error
// 	}

// 	type testcase struct {
// 		expectIn  *api.DeleteReq
// 		expectOut *emptypb.Empty
// 		expectErr error
// 		data      *dataMock
// 	}

// 	testcases := []testcase{
// 		{
// 			// success
// 			expectIn:  &api.DeleteReq{Id: 1},
// 			expectOut: &emptypb.Empty{},
// 			expectErr: nil,
// 			data: &dataMock{
// 				err: nil,
// 			},
// 		},
// 	}

// 	for _, tc := range testcases {
// 		if tc.data == nil {
// 			_repoMock.AssertNotCalled(t, "Delete")
// 		} else {
// 			_repoMock.On("Delete", ctx, tc.expectIn.Id).Return(tc.data.err)
// 		}

// 		actualOut, actualErr := _server.Delete(ctx, tc.expectIn)

// 		assert.Equal(t, tc.expectOut, actualOut)
// 		assert.ErrorIs(t, tc.expectErr, actualErr)
// 	}
// 	_repoMock.AssertExpectations(t)
// }
