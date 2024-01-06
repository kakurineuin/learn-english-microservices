FROM golang:alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN cd ./cmd/examservice && CGO_ENABLED=0 GOOS=linux go build -o /exam-service
EXPOSE 8090 27017
CMD ["/exam-service"]
