# Start with a base Go image
FROM golang:1.22-alpine

# Install SQLite3 dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod tidy

# Copy the source code into the container
COPY . .


# Build the Go app
RUN go build -o conjugador-bot .

# Command to run the executable
CMD ["./conjugador-bot"]
