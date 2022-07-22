package main

import (
	"log"
	"net"
	"strconv"
	"todo-go-grpc/app/dbservice"

	user_service "todo-go-grpc/app/user/api"

	"google.golang.org/grpc"

	tagService "todo-go-grpc/app/task/internal/tag"
	taskService "todo-go-grpc/app/task/internal/task"
	repo "todo-go-grpc/app/task/repository/postgre"
)

const (
	port int = 8082
)

func main() {
	server := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Listen TCP error:\n%v", err)
	}

	db := dbservice.Init()

	taskRepository := repo.NewTaskRepository(*db)
	tagRepository := repo.NewTagRepository(*db)

	// Get user service
	address := "localhost:8081"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connect to user service fail\n%v", err)
	}
	userService := user_service.NewUserHandlerClient(conn)

	taskService.RegisterGrpc(server, taskRepository, tagRepository, userService)
	tagService.RegisterGrpc(server, tagRepository)

	log.Printf("Task service start on port %v", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Unexpected error\n%v", err)
	}
}
