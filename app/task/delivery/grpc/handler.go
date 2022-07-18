package grpc

import (
	"context"
	"errors"
	"log"
	tagDomain "todo-go-grpc/app/tag/domain"
	"todo-go-grpc/app/task/domain"
	_usecase "todo-go-grpc/app/task/usecase"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	usecase _usecase.TaskUsecase
	UnimplementedTaskHandlerServer
}

func NewTagServerGrpc(gserver *grpc.Server, tagUsecase _usecase.TaskUsecase) {
	taskServer := &server{
		usecase: tagUsecase,
	}

	RegisterTaskHandlerServer(gserver, taskServer)
}

func transferDomainToProto(in domain.Task) *Task {
	return &Task{
		Id:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		IsDone:      in.IsDone,
		DonedTime:   timestamppb.New(in.DoneAt),
		CreatorId:   in.CreatorId,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in Task) *domain.Task {
	return &domain.Task{
		ID:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		IsDone:      in.IsDone,
		DoneAt:      in.DonedTime.AsTime(),
		CreatorId:   in.CreatorId,
		CreatedAt:   in.CreatedTime.AsTime(),
	}
}

func (serverInstance *server) List(ctx context.Context, req *ListReq) (*ListTask, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	conditions_map := map[string]any{}
	if req.Name != "" {
		conditions_map["name"] = req.Name
	}
	if req.TagsId != nil || len(req.TagsId) != 0 {
		conditions_map["tags"] = req.TagsId
	}
	if req.Filter != Filter_FILTER_UNSPECIFIED {
		conditions_map["filter"] = req.Filter.String()
	}

	tasks_domain, err := serverInstance.usecase.Fetch(ctx, creator_id, req.PageToken, req.PageSize, conditions_map)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTaskNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	tasks_rs := &ListTask{Tasks: []*Task{}}
	for _, task := range tasks_domain {
		tasks_rs.Tasks = append(tasks_rs.Tasks, transferDomainToProto(task))
	}

	return tasks_rs, nil
}

func (serverInstance *server) Get(ctx context.Context, req *GetReq) (*Task, error) {
	task, err := serverInstance.usecase.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTaskNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*task), nil
}

func (serverInstance *server) Create(ctx context.Context, req *CreateReq) (*Task, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	new_task, err := serverInstance.usecase.Create(ctx, creator_id, &domain.Task{
		Name:        req.Name,
		Description: req.Description,
		IsDone:      req.IsDone,
		Tags:        []tagDomain.Tag{},
	})

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*new_task), nil
}

func (serverInstance *server) Update(ctx context.Context, req *UpdateReq) (*Task, error) {
	data := transferProtoToDomain(*req.NewTaskInfo)
	new_task, err := serverInstance.usecase.Update(ctx, req.Id, data, req.TagsAdded, req.TagsDeleted)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domain.ErrTaskExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*new_task), nil
}

func (serverInstance *server) DeleteMultiple(ctx context.Context, req *DeleteMultipleReq) (*emptypb.Empty, error) {
	err := serverInstance.usecase.Delete(ctx, req.TasksId)

	if err != nil {
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return nil, nil
}

func (serverInstance *server) DeleteAll(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	err := serverInstance.usecase.DeleteAllTaskOfUser(ctx, creator_id)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return nil, nil
}
