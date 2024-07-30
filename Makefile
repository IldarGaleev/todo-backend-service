
OUT_DIR := pkg/grpc/proto


generate: 
	protoc -I proto proto/*.proto \
	--go_out=./$(OUT_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=./$(OUT_DIR) \
	--go-grpc_opt=paths=source_relative