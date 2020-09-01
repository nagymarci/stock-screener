FROM golang:alpine

RUN apk add --no-cache bash

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

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

COPY config.json .

# Export necessary port
EXPOSE 3100

RUN ["chmod", "+x", "/build/wait-for-it.sh"]
# Command to run when starting the container
CMD ["/build/wait-for-it.sh", "stock-scraper:3000", "--", "/dist/main"]