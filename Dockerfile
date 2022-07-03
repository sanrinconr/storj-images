FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update
RUN apk add git

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Run tests
RUN go test ./cmd/...
# Build the application
RUN go build -o main cmd/api/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Command to run when starting the container
CMD ["/dist/main"]

# Build a small image
FROM alpine

RUN apk add --no-cache tzdata
ENV TZ America/Bogota
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY --from=builder /dist/main /

# Command to run
ENTRYPOINT ["/main"]