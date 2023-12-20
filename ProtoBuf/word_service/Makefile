SERVICE_DIR = ../../WordService


build: clean
	protoc --go_out=. --go-grpc_out=. word_service.proto
	mv *.pb.go $(SERVICE_DIR)/pb

clean:
	# Delete all .pb.go files
	find . -type f -name "*.pb.go" -delete
