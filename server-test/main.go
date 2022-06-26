package main

import (
	"github.com/jose827corrza/gogRPC/database"
	"github.com/jose827corrza/gogRPC/server"
	"github.com/jose827corrza/gogRPC/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":5020")
	if err != nil {
		log.Fatal(err)
	}
	repo, err := database.NewPostgresRepository(
		"postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	server := server.NewTestServer(repo)

	s := grpc.NewServer()
	testpb.RegisterTestServiceServer(s, server)
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
