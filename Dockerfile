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
RUN go build -o main .

FROM alpine:latest

RUN apk add --no-cache bash tzdata

# Copy binary from build to main folder
COPY --from=builder /build/main .
COPY wait-for-it.sh .

COPY config.json .

# Export necessary port
EXPOSE 3100

RUN ["chmod", "+x", "./wait-for-it.sh"]
# Command to run when starting the container
CMD ["./wait-for-it.sh", "stock-scraper:3000", "--", "./main"]