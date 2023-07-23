create_all: create_id create_oauth

create_id:
	protoc --proto_path=grpc/proto/id grpc/proto/id/*.proto --go_out=grpc/gen
	protoc --proto_path=grpc/proto/id grpc/proto/id/*.proto --go-grpc_out=grpc/gen

create_oauth:
	protoc --proto_path=grpc/proto/oauth grpc/proto/oauth/*.proto --go_out=grpc/gen
	protoc --proto_path=grpc/proto/oauth grpc/proto/oauth/*.proto --go-grpc_out=grpc/gen
	
clean:
	rm grpc/gen/proto/*/*.go
