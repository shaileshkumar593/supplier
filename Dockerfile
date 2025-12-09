# Stage 1: Build the Go binary
FROM golang:1.25 AS builder
 
ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/go-build-cache
ENV GOMODCACHE=${GOMODCACHE}
ENV GOCACHE=${GOCACHE}
ENV DEBIAN_FRONTEND=noninteractiv
 
# Set the Current Working Directory inside the container
WORKDIR /swallow-supplier
 
RUN which go && go version
# Copy go mod and sum files
COPY go.mod go.sum ./
 
# Download all dependencies
RUN go mod download
 
# Copy the source code into the container
COPY . .
 
# Run go mod tidy and vendor
RUN go mod tidy && go mod vendor
 
# Build the Go app with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
 
# Stage 2: Build the final image
FROM alpine:3.15
 
ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/go-build-cache
ENV GOMODCACHE=${GOMODCACHE}
ENV GOCACHE=${GOCACHE}
ENV DEBIAN_FRONTEND=noninteractive
 
LABEL authors="Shailesh kumar"
 
# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
 
# Install Redis CLI for debugging purposes
RUN apk add --no-cache redis
 
# Set the Current Working Directory inside the container
WORKDIR /swallow-supplier
 
# Copy the pre-built binary file from the previous stage
COPY --from=builder /swallow-supplier/main .
 
# Set environment variable for runtime reference
ENV AUTH_SCOPE_DEF_PATH=/config/auth.yml
 
# ✅ FIX: Use an absolute path in COPY
COPY config/auth.yml ./config/auth.yml
 
COPY excel/YA-GGT-TRIP_Category_Mapping.xlsx /swallow-supplier/excel/YA-GGT-TRIP_Category_Mapping.xlsx
 
# Ensure the executable has the right permissions
RUN chmod +x ./main
 
# Expose port 7001 to the outside world
EXPOSE 7001
 
# Command to run the application
CMD ["./main"]