protoc --proto_path=./proto \
    --go_out=./proto --go_opt=paths=source_relative \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative \
    --validate_out=paths=source_relative,lang=go:./proto  \
   $(find ./proto -name '*.proto')