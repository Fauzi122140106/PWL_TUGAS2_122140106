package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-saga/proto"

	"google.golang.org/grpc"
)

func sagaExecute(orderID string) {
    connOrder, _ := grpc.Dial(":50051", grpc.WithInsecure())
    connPayment, _ := grpc.Dial(":50052", grpc.WithInsecure())
    connShipping, _ := grpc.Dial(":50053", grpc.WithInsecure())

    orderClient := pb.NewOrderServiceClient(connOrder)
    paymentClient := pb.NewPaymentServiceClient(connPayment)
    shippingClient := pb.NewShippingServiceClient(connShipping)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := orderClient.CreateOrder(ctx, &pb.CreateOrderRequest{OrderId: orderID})
    if err != nil {
        log.Println("Order creation failed")
        return
    }

    _, err = paymentClient.ProcessPayment(ctx, &pb.ProcessPaymentRequest{OrderId: orderID})
    if err != nil {
        log.Println("Payment failed, rolling back order")
        orderClient.CancelOrder(ctx, &pb.CancelOrderRequest{OrderId: orderID})
        return
    }

    _, err = shippingClient.StartShipping(ctx, &pb.StartShippingRequest{OrderId: orderID})
    if err != nil {
        log.Println("Shipping failed, rolling back")
        shippingClient.CancelShipping(ctx, &pb.CancelShippingRequest{OrderId: orderID})
        paymentClient.RefundPayment(ctx, &pb.RefundPaymentRequest{OrderId: orderID})
        orderClient.CancelOrder(ctx, &pb.CancelOrderRequest{OrderId: orderID})
        return
    }

    fmt.Println("Transaction successful!")
}

func main() {
    sagaExecute("order_123")
}