
build: clean
	protoc --go_out=. --go-grpc_out=. word_service.proto

clean:
	# Delete all .pb.go files
	find . -type f -name "*.pb.go" -delete
