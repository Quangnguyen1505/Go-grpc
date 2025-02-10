PATH_RUN_SERVER=server/main.go
PATH_RUN_CLIENT=client/main.go

run_server: 
	go run ${PATH_RUN_SERVER}
run_client: 
	go run ${PATH_RUN_CLIENT}

# protoc
protoc_create:
	cd grpc && protoc --go_out=. --go-grpc_out=. order.proto