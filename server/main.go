package server

// Client gRPC Server

import (
	"fmt"
	pb "github.com/vmmgr/controller/proto/proto-go"
	"google.golang.org/grpc"
	"net"
)

const port = ":50200"

type server struct {
	pb.UnimplementedGrpcServer
}

func Server() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGrpcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: %v", err)
	}
}
