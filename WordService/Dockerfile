FROM golang:alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN cd ./cmd/wordservice && CGO_ENABLED=0 GOOS=linux go build -o /word-service
EXPOSE 8091
CMD ["/word-service"]
