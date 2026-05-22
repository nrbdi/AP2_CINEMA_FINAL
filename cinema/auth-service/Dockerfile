FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
ENV PATH="$PATH:/root/go/bin"

# Copy the shared proto folder into the parent directory context
COPY proto /proto

# Copy auth-service specific files from the root build context
COPY auth-service/go.mod ./
COPY auth-service/ .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 50052
CMD ["./auth-service"]