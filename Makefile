.PHONY: gen-category-grpc
gen-grpc:
	protoc -I ./grpc/proto ./grpc/proto/metric.proto \
	--go_out=./grpc/gen/ --go_opt=paths=source_relative \
	--go-grpc_out=./grpc/gen/ --go-grpc_opt=paths=source_relative