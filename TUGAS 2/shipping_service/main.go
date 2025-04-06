package main

import (
	"context"
	"log"
	"net"

	pb "grpc-saga/proto"

	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedShippingServiceServer
}

func (s *server) StartShipping(ctx context.Context, req *pb.StartShippingRequest) (*pb.StartShippingResponse, error) {
    log.Printf("Shipping started for order: %s", req.OrderId)
    return &pb.StartShippingResponse{OrderId: req.OrderId, Status: "SHIPPED"}, nil
}

func (s *server) CancelShipping(ctx context.Context, req *pb.CancelShippingRequest) (*pb.CancelShippingResponse, error) {
    log.Printf("Shipping cancelled for order: %s", req.OrderId)
    return &pb.CancelShippingResponse{Message: "Shipping cancelled"}, nil
}

func main() {
    listener, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterShippingServiceServer(s, &server{})

    log.Println("Shipping Service running on port 50053")
    if err := s.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}