SERVICE_DIR = ../../ExamService


build: clean
	protoc --go_out=. --go-grpc_out=. exam_service.proto
	mv *.pb.go $(SERVICE_DIR)/pb

clean:
	# Delete all .pb.go files
	find . -type f -name "*.pb.go" -delete
