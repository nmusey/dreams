# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Set Go environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create necessary directories with correct permissions
RUN mkdir -p /app/tmp /app/images && \
    chmod -R 777 /app/tmp /app/images

# Create .air.toml for hot reloading configuration
RUN echo '[build]\ncmd = "go build -o ./tmp/main ."\nbin = "./tmp/main"\nfull_bin = "./tmp/main"\ninclude_ext = ["go", "tpl", "tmpl", "html"]\nexclude_dir = ["assets", "tmp", "vendor"]\ninclude_dir = []\nexclude_file = []\ndelay = 1000\nstop_on_error = true\nlog = "air_errors.log"\n\n[log]\ntime = true\n\n[color]\nmain = "magenta"\nwatcher = "cyan"\nbuild = "yellow"\nrunner = "green"\n\n[misc]\nclean_on_exit = true' > .air.toml

# Final stage
FROM golang:1.24-alpine

WORKDIR /app

# Copy the binary and config from builder
COPY --from=builder /app/tmp/main .
COPY --from=builder /app/.air.toml .

# Create necessary directories with correct permissions
RUN mkdir -p /app/tmp /app/images && \
    chmod -R 777 /app/tmp /app/images

# Copy the rest of the application
COPY --from=builder /app .

# Set working directory to /app
WORKDIR /app

# Expose port
EXPOSE 8080

# Use Air for hot reloading
CMD ["air", "-c", ".air.toml"]