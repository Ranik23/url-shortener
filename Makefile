protoc:
	protoc --proto_path=api/proto/ --go_out=. --go-grpc_out=. service.proto

lint:
	golangci-lint run $(filter-out $@, $(MAKECMDGOALS))

.PHONY: lint protoc

