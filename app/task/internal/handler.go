package internal

import (
	"context"
	"errors"
	"log"

	api "todo-go-grpc/app/task/api"
	domain "todo-go-grpc/app/task/domain"
	repository "todo-go-grpc/app/task/repository"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	repo repository.TaskRepository
	api.UnimplementedTaskHandlerServer
}

func RegisterGrpc(gserver *grpc.Server, repo repository.TaskRepository) {
	taskServer := &server{
		repo: repo,
	}

	api.RegisterTaskHandlerServer(gserver, taskServer)
}

func transferDomainToProto(in domain.Task) *api.Task {
	return &api.Task{
		Id:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		IsDone:      in.IsDone,
		DonedTime:   timestamppb.New(in.DoneAt),
		CreatorId:   in.CreatorId,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in api.Task) *domain.Task {
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

func (serverInstance *server) List(ctx context.Context, req *api.ListReq) (*api.ListTask, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	conditions_map := map[string]any{}
	if req.Name != "" {
		conditions_map["name"] = req.Name
	}
	// if req.TagsId != nil || len(req.TagsId) != 0 {
	// 	conditions_map["tags"] = req.TagsId
	// }
	if req.Filter != api.Filter_FILTER_UNSPECIFIED {
		conditions_map["filter"] = req.Filter.String()
	}

	tasks_domain, err := serverInstance.repo.Fetch(ctx, creator_id, req.PageToken, req.PageSize, conditions_map)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTaskNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	tasks_rs := &api.ListTask{Tasks: []*api.Task{}}
	for _, task := range tasks_domain {
		tasks_rs.Tasks = append(tasks_rs.Tasks, transferDomainToProto(task))
	}

	return tasks_rs, nil
}

func (serverInstance *server) Get(ctx context.Context, req *api.GetReq) (*api.Task, error) {
	task, err := serverInstance.repo.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTaskNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*task), nil
}

func (serverInstance *server) Create(ctx context.Context, req *api.CreateReq) (*api.Task, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	data := &domain.Task{
		Name:        req.Name,
		Description: req.Description,
		IsDone:      req.IsDone,
		// Tags:        []tagDomain.Tag{},
	}
	// for _, task_id := range req.Tags {
	// 	data.Tags = append(data.Tags, tagDomain.Tag{ID: task_id})
	// }

	new_task, err := serverInstance.repo.Create(ctx, creator_id, data)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*new_task), nil
}

func (serverInstance *server) Update(ctx context.Context, req *api.UpdateReq) (*api.Task, error) {
	data := transferProtoToDomain(*req.NewTaskInfo)
	new_task, err := serverInstance.repo.Update(ctx, req.Id, data, req.TagsAdded, req.TagsDeleted)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domain.ErrTaskExists) {
			return nil, grpc_status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*new_task), nil
}

func (serverInstance *server) DeleteMultiple(ctx context.Context, req *api.DeleteMultipleReq) (*emptypb.Empty, error) {
	err := serverInstance.repo.Delete(ctx, req.TasksId)

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
	// var creator_id int32 = 1

	err := serverInstance.repo.Delete(ctx, []int32{})

	if err != nil {
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return nil, nil
}
