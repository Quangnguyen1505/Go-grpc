package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	proto "github.com/quangnt/go-grpc/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 9000, "The port to connect")
)

type server struct { //Struct server sẽ implement interface OrderServiceServer
	proto.UnimplementedOrderServiceServer //cung cấp các phương thức mặc định cho interface OrderServiceServer
}

func (s *server) NewOrder(ctx context.Context, in *proto.NewRequestOrder) (*proto.NewResponseOrder, error) {
	log.Printf("Received order:::%v %v %d", in.GetOrderRequest(), in.GetDescription(), in.GetId())

	callbackClient := &proto.NewResponseOrder{
		OrderId: "new orderId " + strconv.Itoa(int(in.GetId())),
		Result:  "success::" + in.GetDescription(),
	}

	product := map[string]interface{}{
		"name":     "áo cổ lọ",
		"color":    "white",
		"quantity": 3,
	}

	go sendOrderToPaymentHandle(strconv.Itoa(int(in.GetId())), product)

	return callbackClient, nil //Trả về response cho client
}

func sendOrderToPaymentHandle(orderId string, product interface{}) {
	address := "localhost:9001"

	// Chuyển đổi product thành struct proto.PaymentDataProduct
	productMap, ok := product.(map[string]interface{})
	if !ok {
		log.Fatal("product is not a valid map[string]interface{}")
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("server order connect to payment error %v", err)
	}
	defer conn.Close()

	c := proto.NewPaymentServiceClient(conn)

	payment := &proto.PaymentRequest{
		OrderId: orderId,
		Product: &proto.PaymentDataProduct{
			Name:     productMap["name"].(string),
			Color:    productMap["color"].(string),
			Quantity: int32(productMap["quantity"].(int)),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.NewPayment(ctx, payment)
	if err != nil {
		log.Fatalf("could not greate: %v", err)
	}
	log.Printf("Status: %s and message: %s", r.GetStatus(), r.GetMessage())
}

func main() {
	fmt.Println("port", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	proto.RegisterOrderServiceServer(s, &server{}) //Đăng ký service OrderServiceServer
	log.Printf("server listening on port %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
