# Build stage
FROM golang:1.24-alpine

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create .air.toml for hot reloading configuration
RUN echo '[build]\ncmd = "go build -o ./tmp/main ."\nbin = "./tmp/main"\nfull_bin = "./tmp/main"\ninclude_ext = ["go", "tpl", "tmpl", "html"]\nexclude_dir = ["assets", "tmp", "vendor"]\ninclude_dir = []\nexclude_file = []\ndelay = 1000\nstop_on_error = true\nlog = "air_errors.log"\n\n[log]\ntime = true\n\n[color]\nmain = "magenta"\nwatcher = "cyan"\nbuild = "yellow"\nrunner = "green"\n\n[misc]\nclean_on_exit = true' > .air.toml

# Expose port
EXPOSE 8080

# Use Air for hot reloading
CMD ["air", "-c", ".air.toml"] 