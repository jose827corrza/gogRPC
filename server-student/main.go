package main

import (
	"github.com/jose827corrza/gogRPC/database"
	"github.com/jose827corrza/gogRPC/server"
	"github.com/jose827corrza/gogRPC/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":5010")
	if err != nil {
		log.Fatal(err)
	}
	repo, err := database.NewPostgresRepository(
		"postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	server := server.NewServer(repo)

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, server)
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
