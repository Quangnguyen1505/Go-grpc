package main

import (
	"context"
	"log"
	"time"

	proto "github.com/quangnt/go-grpc/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:9000"
	// Kết nối đến server gRPC (không cần chứng chỉ TLS).
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("didn't connect: %v ", err)
	}
	defer conn.Close()
	// Tạo một client stub để gọi các phương thức được định nghĩa trong OrderService.
	c := proto.NewOrderServiceClient(conn)

	ticker := time.NewTicker(2 * time.Second) // trong trường hợp này ticket được khởi tạo với thời gian 2 giây
	defer ticker.Stop()

	for range ticker.C { // la mot cach tuyet voi de lặp vô hạn và chờ đợi tín hiệu trong go
		orderId := "1001" // make([]byte, 10*1024*1024) //10mb
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		r, err := c.NewOrder(ctx, &proto.NewRequestOrder{OrderRequest: string(orderId)})
		if err != nil {
			log.Fatalf("could not greate: %v", err)
		}
		log.Printf("Order: %s", r.GetOrderResponse())
		cancel()
	}
}
