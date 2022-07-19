package main

import (
	"log"
	"net"
	"strconv"
	"todo-go-grpc/app/dbservice"

	"google.golang.org/grpc"

	service "todo-go-grpc/app/user/internal"
	repo "todo-go-grpc/app/user/repository/postgre"
)

const (
	port int = 8081
)

func main() {
	server := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Listen TCP error: %v", err)
	}

	db := dbservice.Init()

	userRepository := repo.NewUserRepository(*db)
	service.RegisterGrpc(server, userRepository)

	log.Printf("User service start on port %v", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Unexpected error\n%v", err)
	}
}
