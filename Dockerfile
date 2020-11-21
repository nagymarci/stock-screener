FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./cmd/server

FROM alpine:latest

RUN apk add --no-cache bash tzdata

# Copy binary from build to main folder
COPY --from=builder /build/main .

# Export necessary port
EXPOSE 3100

# Command to run when starting the container
CMD ["./main"]