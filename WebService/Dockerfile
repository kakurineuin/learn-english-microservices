FROM golang:alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN apk add --update npm
RUN cd ./frontend && npm install && npm run build
RUN cd ./cmd/webservice && CGO_ENABLED=0 GOOS=linux go build -o /web-service
EXPOSE 8080
CMD ["/web-service"]
