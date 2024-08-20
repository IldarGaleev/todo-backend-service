

ifeq ($(OS),Windows_NT)
RM_TOOL:=del
else
RM_TOOL:=rm
endif


OUT_DIR := pkg/grpc/proto

PHONY: generate cover

generate: 
	protoc -I proto proto/*.proto \
	--go_out=./$(OUT_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=./$(OUT_DIR) \
	--go-grpc_opt=paths=source_relative


cover:
	go test -short -count=1 -coverprofile=.\coverage.out ./...
	go tool cover -html=.\coverage.out
	@$(RM_TOOL) .\coverage.out