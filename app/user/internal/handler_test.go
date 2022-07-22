package internal

import (
	"context"
	"testing"
	"time"
	"todo-go-grpc/app/user/api"
	"todo-go-grpc/app/user/domain"
	"todo-go-grpc/app/user/repository/mocks"

	response_service "todo-go-grpc/app/response_handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestLogin(t *testing.T) {
	_repoMock := &mocks.UserRepository{}
	_server := server{
		repo: _repoMock,
	}
	ctx := context.Background()

	type dataMock struct {
		user *domain.User
		err  error
	}

	type testcase struct {
		expectIn  *api.LoginReq
		expectOut *api.BasicUser
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn:  &api.LoginReq{Username: "vungocthach", Password: "123456"},
			expectOut: &api.BasicUser{Id: 1, Name: "Vũ Ngọc Thạch", Username: "vungocthach", Password: "123456"},
			expectErr: nil,
			data: &dataMock{
				user: &domain.User{ID: 1, Name: "Vũ Ngọc Thạch", Username: "vungocthach", Password: "123456", CreatedAt: time.Now()},
				err:  nil,
			},
		},
		{
			// fail: wrong username or password
			expectIn:  &api.LoginReq{Username: "vungocthach", Password: "12345678"},
			expectOut: nil,
			expectErr: response_service.ResponseErrorNotFound(domain.ErrUserNotExists),
			data: &dataMock{
				user: nil,
				err:  domain.ErrUserNotExists,
			},
		},
		{
			// fail: invalid username
			expectIn:  &api.LoginReq{Username: "", Password: "123456"},
			expectOut: nil,
			expectErr: response_service.ResponseErrorInvalidArgument(api.ErrUsernameOrPaswordIsEmpty),
			data:      nil,
		},
		{
			// fail: invalid password
			expectIn:  &api.LoginReq{Username: "vungocthach", Password: ""},
			expectOut: nil,
			expectErr: response_service.ResponseErrorInvalidArgument(api.ErrUsernameOrPaswordIsEmpty),
			data:      nil,
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_repoMock.AssertNotCalled(t, "GetByUsernameAndPassword")
		} else {
			_repoMock.On("GetByUsernameAndPassword", ctx, tc.expectIn.Username, tc.expectIn.Password).Return(tc.data.user, tc.data.err)
		}

		actualOut, actualErr := _server.Login(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_repoMock.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	_repoMock := &mocks.UserRepository{}
	_server := server{
		repo: _repoMock,
	}
	ctx := context.Background()

	type dataMock struct {
		user *domain.User
		err  error
	}

	type testcase struct {
		expectIn  *api.CreateReq
		expectOut *api.User
		expectErr error
		data      *dataMock
	}

	created_time := time.Now()
	testcases := []testcase{
		{
			// success
			expectIn:  &api.CreateReq{Name: "Vũ Ngọc Thạch", Username: "vungocthach", Password: "123456"},
			expectOut: &api.User{Id: 1, Name: "Vũ Ngọc Thạch", Username: "vungocthach", Password: "123456", CreatedTime: timestamppb.New(created_time)},
			expectErr: nil,
			data: &dataMock{
				user: &domain.User{ID: 1, Name: "Vũ Ngọc Thạch", Username: "vungocthach", Password: "123456", CreatedAt: created_time},
				err:  nil,
			},
		},
		{
			// fail: duplicate username
			expectIn:  &api.CreateReq{Name: "Minh Nhực", Username: "minhnhuc", Password: "123456"},
			expectOut: nil,
			expectErr: response_service.ResponseErrorAlreadyExists(domain.ErrUserNameIsExists),
			data: &dataMock{
				user: nil,
				err:  domain.ErrUserNameIsExists,
			},
		},
		{
			// fail: invalid name
			expectIn:  &api.CreateReq{Name: "", Username: "vungocthach", Password: "123456"},
			expectOut: nil,
			expectErr: response_service.ResponseErrorInvalidArgument(api.ErrCreateReqIsEmpty),
			data:      nil,
		},
		{
			// fail: invalid username
			expectIn:  &api.CreateReq{Name: "Vũ Ngọc Thạch", Username: "", Password: "123456"},
			expectOut: nil,
			expectErr: response_service.ResponseErrorInvalidArgument(api.ErrCreateReqIsEmpty),
			data:      nil,
		},
		{
			// fail: invalid password
			expectIn:  &api.CreateReq{Name: "Vũ Ngọc Thạch", Username: "vungocthach", Password: ""},
			expectOut: nil,
			expectErr: response_service.ResponseErrorInvalidArgument(api.ErrCreateReqIsEmpty),
			data:      nil,
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_repoMock.AssertNotCalled(t, "Create")
		} else {
			_repoMock.On("Create", ctx, &domain.User{
				Name:     tc.expectIn.Name,
				Username: tc.expectIn.Username,
				Password: tc.expectIn.Password,
			},
			).Return(tc.data.user, tc.data.err)
		}

		actualOut, actualErr := _server.Create(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_repoMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	_repoMock := &mocks.UserRepository{}
	_server := server{
		repo: _repoMock,
	}
	ctx := context.Background()

	type dataMock struct {
		user *domain.User
		err  error
	}

	type testcase struct {
		expectIn  *api.UpdateReq
		expectOut *api.User
		expectErr error
		data      *dataMock
	}

	created_time := time.Now()
	testcases := []testcase{
		{
			// success
			expectIn: &api.UpdateReq{
				Id: 1,
				NewUserInfor: &api.User{
					Name:        "Vũ Ngọc Thạch",
					Username:    "vungocthach",
					Password:    "123456",
					CreatedTime: timestamppb.New(created_time),
				},
			},
			expectOut: &api.User{
				Id:          1,
				Name:        "Vũ Ngọc Thạch",
				Username:    "vungocthach",
				Password:    "123456",
				CreatedTime: timestamppb.New(created_time),
			},
			expectErr: nil,
			data: &dataMock{
				user: &domain.User{
					ID:        1,
					Name:      "Vũ Ngọc Thạch",
					Username:  "vungocthach",
					Password:  "123456",
					CreatedAt: created_time,
				},
				err: nil,
			},
		},
		{
			// fail: duplicate username
			expectIn: &api.UpdateReq{
				Id: 1,
				NewUserInfor: &api.User{
					Name:     "Minh Nhực",
					Username: "minhnhuc",
					Password: "123456",
				},
			},
			expectOut: nil,
			expectErr: response_service.ResponseErrorAlreadyExists(domain.ErrUserNameIsExists),
			data: &dataMock{
				user: nil,
				err:  domain.ErrUserNameIsExists,
			},
		},
		{
			// fail: id not exists
			expectIn: &api.UpdateReq{
				Id:           2,
				NewUserInfor: &api.User{},
			},
			expectOut: nil,
			expectErr: response_service.ResponseErrorNotFound(domain.ErrUserNotExists),
			data:      nil,
		},
	}

	// If id = 1, found user, else if = 2, not found
	_repoMock.On("GetByID", ctx, mock.MatchedBy(func(id int32) bool { return id == 1 })).Return(nil, nil)
	_repoMock.On("GetByID", ctx, mock.MatchedBy(func(id int32) bool { return id == 2 })).Return(nil, domain.ErrUserNotExists)

	for _, tc := range testcases {
		if tc.data == nil {
			_repoMock.AssertNotCalled(t, "Update")
		} else {
			_repoMock.On("Update", ctx, tc.expectIn.Id,
				&domain.User{
					// ID:        tc.expectIn.Id,
					Name:      tc.expectIn.NewUserInfor.Name,
					Username:  tc.expectIn.NewUserInfor.Username,
					Password:  tc.expectIn.NewUserInfor.Password,
					CreatedAt: tc.expectIn.NewUserInfor.CreatedTime.AsTime(),
				},
			).Return(tc.data.user, tc.data.err)
		}

		actualOut, actualErr := _server.Update(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_repoMock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	_repoMock := &mocks.UserRepository{}
	_server := server{
		repo: _repoMock,
	}
	ctx := context.Background()

	type dataMock struct {
		err error
	}

	type testcase struct {
		expectIn  *api.DeleteReq
		expectOut *emptypb.Empty
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn:  &api.DeleteReq{Id: 1},
			expectOut: &emptypb.Empty{},
			expectErr: nil,
			data: &dataMock{
				err: nil,
			},
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_repoMock.AssertNotCalled(t, "Delete")
		} else {
			_repoMock.On("Delete", ctx, tc.expectIn.Id).Return(tc.data.err)
		}

		actualOut, actualErr := _server.Delete(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_repoMock.AssertExpectations(t)
}
