package main

import (
	"context"
	"log"
	"net"

	pb "grpc-saga/proto"

	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedOrderServiceServer
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    log.Printf("Order created: %s", req.OrderId)
    return &pb.CreateOrderResponse{OrderId: req.OrderId, Status: "PENDING"}, nil
}

func (s *server) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
    log.Printf("Order cancelled: %s", req.OrderId)
    return &pb.CancelOrderResponse{Message: "Order cancelled"}, nil
}

func main() {
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterOrderServiceServer(s, &server{})

    log.Println("Order Service running on port 50051")
    if err := s.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}