SERVICE_DIR = ../../WordService
WEB_SERVICE_DIR = ../../WebService


build: clean
	protoc --go_out=. --go-grpc_out=. word_service.proto
	cp *.pb.go $(SERVICE_DIR)/pb
	cp *.pb.go $(WEB_SERVICE_DIR)/pb
	rm *.pb.go

clean:
	# Delete all .pb.go files
	find . -type f -name "*.pb.go" -delete
