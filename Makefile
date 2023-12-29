path_proto:=./pb
proto:
	protoc --go_out=${path_proto} --go-grpc_out=${path_proto} ${path_proto}/*.proto

tidy:
	go mod tidy