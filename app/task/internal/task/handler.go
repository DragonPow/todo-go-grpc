package internal

import (
	"context"
	"errors"
	"log"
	"sync"

	response_service "todo-go-grpc/app/response_handler"
	api "todo-go-grpc/app/task/api/task"
	domain "todo-go-grpc/app/task/domain"
	repository "todo-go-grpc/app/task/repository"
	user_service "todo-go-grpc/app/user/api"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	taskRepo repository.TaskRepository
	tagRepo  repository.TagRepository
	api.UnimplementedTaskHandlerServer
	userService user_service.UserHandlerClient
}

func RegisterGrpc(gserver *grpc.Server, taskRepo repository.TaskRepository, tagRepo repository.TagRepository, userService user_service.UserHandlerClient) {
	taskServer := &server{
		taskRepo:    taskRepo,
		tagRepo:     tagRepo,
		userService: userService,
	}

	api.RegisterTaskHandlerServer(gserver, taskServer)
}

// func getUserService() (user_service.UserHandlerClient, error) {
// 	address := "localhost:8081"
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user_service.NewUserHandlerClient(conn), nil
// }

func transferDomainToTag(in *domain.Tag) *api.Tag {
	return &api.Tag{
		Id:          in.ID,
		Value:       in.Value,
		Description: in.Description,
	}
}

func transferTagToDomain(in *api.Tag) *domain.Tag {
	return &domain.Tag{
		ID:          in.Id,
		Value:       in.Value,
		Description: in.Description,
	}
}

func transferDomainToTask(in *domain.Task, creator *api.User) *api.Task {
	apiTasks := []*api.Tag{}
	for _, tag := range in.Tags {
		apiTasks = append(apiTasks, transferDomainToTag(&tag))
	}
	return &api.Task{
		Id:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		IsDone:      in.IsDone,
		Tags:        apiTasks,
		Creator:     creator,
		DonedTime:   timestamppb.New(in.DoneAt),
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferTaskToDomain(in *api.Task) *domain.Task {
	domainTasks := []domain.Tag{}
	for _, tag := range in.Tags {
		domainTasks = append(domainTasks, *transferTagToDomain(tag))
	}
	return &domain.Task{
		ID:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		IsDone:      in.IsDone,
		DoneAt:      in.DonedTime.AsTime(),
		Tags:        domainTasks,
		CreatedAt:   in.CreatedTime.AsTime(),
	}
}

func transferDomainToBasicTask(in *domain.Task) *api.BasicTask {
	rs := &api.BasicTask{
		Id:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		IsDone:      in.IsDone,
		DonedTime:   timestamppb.New(in.DoneAt),
		CreatorId:   in.CreatorId,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
	rs.TagsId = []int32{}
	for _, tag := range in.Tags {
		rs.TagsId = append(rs.TagsId, tag.ID)
	}
	return rs
}

func transferBasicTaskToDomain(in *api.BasicTask) *domain.Task {
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

func (serverInstance *server) GetUserInfo(ctx context.Context, id int32) (*api.User, error) {
	user, err := serverInstance.userService.Get(ctx, &user_service.GetReq{Id: id})
	if err != nil {
		return nil, err
	}

	apiUser := &api.User{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
	}

	return apiUser, nil
}

func (serverInstance *server) List(ctx context.Context, req *api.ListReq) (*api.ListTask, error) {
	// TODO: Get creator id
	var creator_id int32 = 1
	var wg sync.WaitGroup

	// Map req data to search conditions
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

	tasks_domain, err := serverInstance.taskRepo.Fetch(ctx, creator_id, req.PageToken, req.PageSize, conditions_map)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTaskNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	// Tranfer domain to api response
	tasks_rs := &api.ListTask{Tasks: []*api.Task{}}
	for _, task := range tasks_domain {
		my_task := task
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Get user info
			apiUser, _ := serverInstance.GetUserInfo(ctx, my_task.CreatorId)
			// TODO: handler error here
			// if err != nil {
			// 	if errors.Is(err, domain.ErrUserNotExists) {
			// 		panic(grpc_status.Error(codes.NotFound, err.Error()))
			// 	}
			// 	panic(grpc_status.Error(codes.Unknown, err.Error()))
			// }

			tasks_rs.Tasks = append(tasks_rs.Tasks, transferDomainToTask(&my_task, apiUser))
		}()
	}
	wg.Wait()

	return tasks_rs, nil
}

func (serverInstance *server) Get(ctx context.Context, req *api.GetReq) (*api.Task, error) {
	// Get task
	task, err := serverInstance.taskRepo.GetByID(ctx, req.Id)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTaskNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	// Get user info
	apiUser, err := serverInstance.GetUserInfo(ctx, task.CreatorId)
	if err != nil {
		log.Fatalf("Error getting user info: %v", err)
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}
	log.Printf("Get user info success, id: %v", apiUser.Id)

	return transferDomainToTask(task, apiUser), nil
}

func (serverInstance *server) Create(ctx context.Context, req *api.CreateReq) (*api.BasicTask, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	// Check user not found
	if _, err := serverInstance.GetUserInfo(ctx, creator_id); err != nil {
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	// Tranfer req data to domain
	data := &domain.Task{
		Name:        req.Name,
		Description: req.Description,
		IsDone:      req.IsDone,
		Tags:        []domain.Tag{},
	}
	for _, task_id := range req.Tags {
		data.Tags = append(data.Tags, domain.Tag{ID: task_id})
	}

	new_task, err := serverInstance.taskRepo.Create(ctx, creator_id, data)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToBasicTask(new_task), nil
}

func (serverInstance *server) Update(ctx context.Context, req *api.UpdateReq) (*api.BasicTask, error) {
	data := transferBasicTaskToDomain(req.NewTaskInfo)

	// Check tag add and remove is in tags
	all_tag, err := serverInstance.tagRepo.FetchAll(ctx)
	if err != nil {
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	// Check add_tag in tags
	for _, add_tag := range req.TagsAdded {
		isFind := false
		for _, tag := range all_tag {
			if add_tag == tag.ID {
				isFind = true
				break
			}
		}
		if !isFind {
			return nil, grpc_status.Error(codes.NotFound, domain.ErrTagNotExists.Error())
		}
	}

	// Check delete_tag in tags
	for _, delete_tag := range req.TagsDeleted {
		isFind := false
		for _, tag := range all_tag {
			if delete_tag == tag.ID {
				isFind = true
				break
			}
		}
		if !isFind {
			return nil, grpc_status.Error(codes.NotFound, domain.ErrTagNotExists.Error())
		}
	}

	new_task, err := serverInstance.taskRepo.Update(ctx, req.Id, data, req.TagsAdded, req.TagsDeleted)
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

	return transferDomainToBasicTask(new_task), nil
}

func (serverInstance *server) DeleteMultiple(ctx context.Context, req *api.DeleteMultipleReq) (*emptypb.Empty, error) {
	if err := serverInstance.taskRepo.Delete(ctx, req.TasksId); err != nil {
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (serverInstance *server) DeleteAll(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO: Get creator id
	var creator_id int32 = 1

	tasks_id, err := serverInstance.taskRepo.GetByUserId(ctx, creator_id)
	if err != nil {
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	err = serverInstance.taskRepo.Delete(ctx, tasks_id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}
