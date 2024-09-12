proto:
	mkdir -p pkg/proto/v1/grpc
	protoc --go_out=./pkg/proto/v1/grpc --go_opt=paths=source_relative \
		--proto_path=pkg/proto/v1 \
		--go-grpc_out=./pkg/proto/v1/grpc --go-grpc_opt=paths=source_relative \
		pkg/proto/v1/*.proto

install-dependencies:
	sudo sh ./scripts/install-dependencies.sh pn532_uart 0