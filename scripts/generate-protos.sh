#!bash

protoc --go_out=internal/ --proto_path=proto/ proto/grpc_msg.proto
protoc --go-grpc_out=internal/ --proto_path=proto/ proto/grpc_services.proto
