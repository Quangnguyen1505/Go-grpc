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
		orderId := int32(1001) // make([]byte, 10*1024*1024) //10mb
		description := "total bill buy pizza"
		// Gộp thành JSON
		// dataResult := map[string]interface{}{
		// 	"orderId":     orderId,
		// 	"description": description,
		// }

		// fmt.Println("dataResult::", dataResult)

		// // Chuyển map thành JSON (marshal)
		// jsonData, err := json.Marshal(dataResult)
		// if err != nil {
		// 	log.Fatalf("Error marshaling JSON: %v", err)
		// }

		// // in ra chuoi byte
		// fmt.Println("jsonData::", jsonData)
		// // In ra chuỗi JSON thay vì mảng số
		// log.Println("JSON string:", string(jsonData))

		order := &proto.NewRequestOrder{
			OrderRequest: "bill bill",
			Description:  description,
			Id:           orderId,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		r, err := c.NewOrder(ctx, order)
		if err != nil {
			log.Fatalf("could not greate: %v", err)
		}
		log.Printf("OrderId: %s and result: %s", r.GetOrderId(), r.GetResult())
		cancel()
	}
}
