FROM golang:alpine3.19

# Redis
RUN apk add --update redis
WORKDIR /app

# ExamService
WORKDIR ExamService
COPY ./ExamService/go.mod ./ExamService/go.sum ./
RUN go mod download
COPY ./ExamService ./
RUN cd ./cmd/examservice && CGO_ENABLED=0 GOOS=linux go build -o /exam-service
EXPOSE 8090

# WordService
WORKDIR ../WordService
COPY ./WordService/go.mod ./WordService/go.sum ./
RUN go mod download
COPY ./WordService ./
RUN cd ./cmd/wordservice && CGO_ENABLED=0 GOOS=linux go build -o /word-service
EXPOSE 8091

# WebService
WORKDIR ../WebService
COPY ./WebService/go.mod ./WebService/go.sum ./
RUN go mod download
COPY ./WebService ./
RUN apk add --update npm
RUN cd ./frontend && npm install && npm run build
RUN cd ./cmd/webservice && CGO_ENABLED=0 GOOS=linux go build -o /web-service
EXPOSE 8080
CMD ["sh", "-c", "redis-server & /exam-service & /word-service & sleep 5; /web-service"]
