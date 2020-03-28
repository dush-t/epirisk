# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Dushyant Yadav <dushyant9309@gmail.com>"

# Add cgo dependencies
RUN apk add --no-cache ca-certificates cmake make g++ openssl-dev git curl pkgconfig

# Install seabolt
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt 
WORKDIR /seabolt/build
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 