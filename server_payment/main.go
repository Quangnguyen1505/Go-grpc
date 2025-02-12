package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	proto "github.com/quangnt/go-grpc/grpc"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9001, "The port to connect")
)

type server struct {
	proto.UnimplementedPaymentServiceServer
}

func (s *server) NewPayment(ctx context.Context, in *proto.PaymentRequest) (*proto.PaymentResponse, error) {
	log.Printf("Reveived for order %s with %v\n", in.GetOrderId(), in.GetProduct())
	callbackOrder := &proto.PaymentResponse{
		Status:  200,
		Message: "payment successs",
	}

	fmt.Println("Process handle .....")

	// time.Sleep(2 * time.Second)
	return callbackOrder, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("error listen server payment %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterPaymentServiceServer(s, &server{})
	log.Printf("server listening on port %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
