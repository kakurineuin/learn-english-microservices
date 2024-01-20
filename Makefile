EXAM_SERVICE_DIR = ./ExamService
WORD_SERVICE_DIR = ./WordService
WEB_SERVICE_DIR = ./WebService

test: build
	# Test ExamService
	cd $(EXAM_SERVICE_DIR) && go generate ./... && go clean -testcache && go test ./...

	# Test WordService
	cd $(WORD_SERVICE_DIR) && go generate ./... && go clean -testcache && go test ./...

	# Test WebService
	cd $(WEB_SERVICE_DIR) && go generate ./... && go clean -testcache && go test ./...

build:
	docker build $(EXAM_SERVICE_DIR) -t mises/exam-service:test
	docker build $(WORD_SERVICE_DIR) -t mises/word-service:test

