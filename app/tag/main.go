package main

import (
	"log"
	"net"
	"strconv"
	"todo-go-grpc/app/dbservice"

	"google.golang.org/grpc"

	service "todo-go-grpc/app/tag/internal"
	repo "todo-go-grpc/app/tag/repository/postgre"
)

const (
	port int = 8083
)

func main() {
	server := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Listen TCP error:\n%v", err)
	}

	db := dbservice.Init()

	tagRepository := repo.NewTagRepository(*db)
	service.RegisterGrpc(server, tagRepository)

	log.Printf("Tag service start on port %v", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Unexpected error\n%v", err)
	}
}
