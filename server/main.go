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
	port = flag.Int("port", 9000, "The port to connect")
)

type server struct { //Struct server sẽ implement interface OrderServiceServer
	proto.UnimplementedOrderServiceServer //cung cấp các phương thức mặc định cho interface OrderServiceServer
}

func (s *server) NewOrder(ctx context.Context, in *proto.NewRequestOrder) (*proto.NewResponseOrder, error) {
	log.Printf("Received order:::%v", in.GetOrderRequest())
	return &proto.NewResponseOrder{OrderResponse: "new orderId " + in.GetOrderRequest()}, nil
}

func main() {
	fmt.Println("port", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	proto.RegisterOrderServiceServer(s, &server{})
	log.Printf("server listening on port %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
