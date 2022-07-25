package internal

import (
	"context"
	"log"
	"testing"
	"time"
	"todo-go-grpc/app/response_handler"
	api "todo-go-grpc/app/task/api/tag"
	"todo-go-grpc/app/task/domain"
	"todo-go-grpc/app/task/repository/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	time_now     = time.Now()
	time_init    = time.Unix(0, 0).UTC()
	_tagRepoMock = &mocks.TagRepository{}
	_server      = server{
		repo: _tagRepoMock,
	}
	ctx = context.Background()
)

func InitMock(t *testing.T) {

}

func TestList(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		tag []domain.Tag
		err error
	}

	type testcase struct {
		expectIn  *api.ListReq
		expectOut *api.ListTag
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn: &api.ListReq{},
			expectOut: &api.ListTag{
				Tags: []*api.Tag{
					{
						Id:          1,
						Value:       "Tag 1",
						Description: "No description",
						CreatedTime: timestamppb.New(time_now),
					},
					{
						Id:          2,
						Value:       "Tag 2",
						Description: "No description",
						CreatedTime: timestamppb.New(time_now),
					},
				},
			},
			data: &dataMock{
				tag: []domain.Tag{
					{
						ID:          1,
						Value:       "Tag 1",
						Description: "No description",
						CreatedAt:   time_now,
					},
					{
						ID:          2,
						Value:       "Tag 2",
						Description: "No description",
						CreatedAt:   time_now,
					},
				},
				err: nil,
			},
		},
	}

	for index, tc := range testcases {
		if tc.data == nil {
			_tagRepoMock.AssertNotCalled(t, "FetchAll")
		} else {
			_tagRepoMock.On("FetchAll", ctx).Return(tc.data.tag, tc.data.err)
		}

		actualOut, actualErr := _server.List(ctx, tc.expectIn)

		log.Printf("Output assert, testcase %v: %v", index+1, assert.Equal(t, tc.expectOut, actualOut))
		log.Printf("Error assert, testcase %v: %v", index+1, assert.ErrorIs(t, tc.expectErr, actualErr))
	}
	_tagRepoMock.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		tag *domain.Tag
		err error
	}

	type testcase struct {
		expectIn  *api.GetReq
		expectOut *api.Tag
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn: &api.GetReq{Id: 1},
			expectOut: &api.Tag{
				Id:          1,
				Value:       "Tag 1",
				Description: "No description",
				CreatedTime: timestamppb.New(time_now),
			},
			data: &dataMock{
				tag: &domain.Tag{
					ID:          1,
					Value:       "Tag 1",
					Description: "No description",
					CreatedAt:   time_now,
				},
				err: nil,
			},
		},
		{
			// success
			expectIn:  &api.GetReq{Id: 2},
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrTagNotExists),
			data: &dataMock{
				tag: nil,
				err: domain.ErrTagNotExists,
			},
		},
	}

	for index, tc := range testcases {
		if tc.data == nil {
			_tagRepoMock.AssertNotCalled(t, "GetByID")
		} else {
			_tagRepoMock.On("GetByID", ctx, tc.expectIn.Id).Return(tc.data.tag, tc.data.err)
		}

		actualOut, actualErr := _server.Get(ctx, tc.expectIn)

		log.Printf("Output assert, testcase %v: %v", index+1, assert.Equal(t, tc.expectOut, actualOut))
		log.Printf("Error assert, testcase %v: %v", index+1, assert.ErrorIs(t, tc.expectErr, actualErr))
	}
	_tagRepoMock.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	InitMock(t)
	// define testcase
	type dataMock struct {
		tag *domain.Tag
		err error
	}

	type testcase struct {
		expectIn  *api.CreateReq
		expectOut *api.Tag
		expectErr error
		data      *dataMock
	}

	// create testcase
	testcases := []testcase{
		{
			// success
			expectIn: &api.CreateReq{
				Value:       "Tag value",
				Description: "No comment",
			},
			expectOut: &api.Tag{
				Id:          1,
				Value:       "Tag value",
				Description: "No comment",
				CreatedTime: timestamppb.New(time_now),
			},
			expectErr: nil,
			data: &dataMock{
				tag: &domain.Tag{
					ID:          1,
					Value:       "Tag value",
					Description: "No comment",
					CreatedAt:   time_now,
				},
			},
		},
		{
			// fail: tag is exists
			expectIn: &api.CreateReq{
				Value:       "Tag duplicate",
				Description: "No comment",
			},
			expectErr: response_handler.ResponseErrorAlreadyExists(domain.ErrTagIsExists),
			data: &dataMock{
				tag: nil,
				err: domain.ErrTagIsExists,
			},
		},
	}

	// test
	for _, tc := range testcases {
		if tc.data == nil {
			_tagRepoMock.AssertNotCalled(t, "Create")
		} else {
			// Tranfer data
			dataTransfer := &domain.Tag{
				Value:       tc.expectIn.Value,
				Description: tc.expectIn.Description,
			}

			_tagRepoMock.On("Create", ctx, dataTransfer).Return(tc.data.tag, tc.data.err)
		}

		actualOut, actualErr := _server.Create(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_tagRepoMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	InitMock(t)

	type dataMock struct {
		task *domain.Tag
		err  error
	}

	type testcase struct {
		expectIn  *api.UpdateReq
		expectOut *api.Tag
		expectErr error
		data      *dataMock
	}

	testcases := []testcase{
		{
			// success
			expectIn: &api.UpdateReq{
				Id:         1,
				NewTagInfo: &api.Tag{Value: "Tag update", Description: "No description"},
			},
			expectOut: &api.Tag{
				Id:          1,
				Value:       "Tag update",
				Description: "No description",
				CreatedTime: timestamppb.New(time_now),
			},
			data: &dataMock{
				task: &domain.Tag{
					ID:          1,
					Value:       "Tag update",
					Description: "No description",
					CreatedAt:   time_now,
				},
				err: nil,
			},
		},
		{
			// fail: tag is exists
			expectIn: &api.UpdateReq{
				Id:         2,
				NewTagInfo: &api.Tag{Value: "Tag exists", Description: "No description"},
			},
			expectErr: response_handler.ResponseErrorAlreadyExists(domain.ErrTagIsExists),
			data: &dataMock{
				task: nil,
				err:  domain.ErrTagIsExists,
			},
		},
		{
			// fail: tag not found
			expectIn: &api.UpdateReq{
				Id:         3,
				NewTagInfo: &api.Tag{Value: "Tag not found", Description: "No description"},
			},
			expectErr: response_handler.ResponseErrorNotFound(domain.ErrTagNotExists),
			data: &dataMock{
				task: nil,
				err:  domain.ErrTagNotExists,
			},
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_tagRepoMock.AssertNotCalled(t, "Update")
		} else {
			data := transferProtoToDomain(tc.expectIn.NewTagInfo)
			_tagRepoMock.On("Update", ctx, tc.expectIn.Id, data).Return(tc.data.task, tc.data.err)
		}

		actualOut, actualErr := _server.Update(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_tagRepoMock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	InitMock(t)

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
			expectIn: &api.DeleteReq{
				Id: 1,
			},
			expectOut: &emptypb.Empty{},
			data: &dataMock{
				err: nil,
			},
		},
	}

	for _, tc := range testcases {
		if tc.data == nil {
			_tagRepoMock.AssertNotCalled(t, "Delete")
		} else {
			_tagRepoMock.On("Delete", ctx, tc.expectIn.Id).Return(tc.data.err)
		}

		actualOut, actualErr := _server.Delete(ctx, tc.expectIn)

		assert.Equal(t, tc.expectOut, actualOut)
		assert.ErrorIs(t, tc.expectErr, actualErr)
	}
	_tagRepoMock.AssertExpectations(t)
}
