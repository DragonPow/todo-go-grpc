package main

import (
	"log"
	"net"

	"todo-go-grpc/app/dbservice"

	"google.golang.org/grpc"

	// _tagRepository "todo-go-grpc/app/tag/repository/postgre"
	// _taskRepository "todo-go-grpc/app/task/repository/postgre"
	_userRepository "todo-go-grpc/app/user/repository/postgre"

	// _tagUsecase "todo-go-grpc/app/tag/usecase"
	// _taskUsecase "todo-go-grpc/app/task/usecase"
	_userUsecase "todo-go-grpc/app/user/usecase"

	// _tagService "todo-go-grpc/app/tag/delivery/grpc"
	// _taskService "todo-go-grpc/app/task/delivery/grpc"
	_userService "todo-go-grpc/app/user/delivery/grpc"
)

func main() {
	server := grpc.NewServer()

	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatalf("Listen TCP error:\n%v", err)
	}

	db := dbservice.Init()

	// tagRepository := _tagRepository.NewTagRepository(*db)
	// taskRepository := _taskRepository.NewTaskRepository(*db)
	userRepository := _userRepository.NewUserRepository(*db)

	// tagUsecase := _tagUsecase.NewTagUsecase(*db, tagRepository)
	// taskUsecase := _taskUsecase.NewTaskUsecase(*db, taskRepository, userRepository, tagRepository)
	userUsecase := _userUsecase.NewUserUsecase(*db, userRepository)

	// _tagService.RegisterGrpc(server, tagUsecase)
	// _taskService.RegisterGrpc(server, taskUsecase)
	_userService.RegisterGrpc(server, userUsecase)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Unexpected error", err)
	}
}
