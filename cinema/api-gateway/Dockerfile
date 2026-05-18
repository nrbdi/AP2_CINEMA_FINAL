FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
ENV PATH="$PATH:/root/go/bin"


COPY proto /proto


COPY api-gateway/go.mod ./
COPY api-gateway/ .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]