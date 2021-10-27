compile:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	proto/hello.proto

build_client:
	docker build -t primozh/grpc-go-test-client -f Dockerfile.client .

build_server:
	docker build -t primozh/grpc-go-test-server -f Dockerfile.server .
