FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN mkdir -p /app/tmp /app/images && \
    chmod -R 777 /app/tmp /app/images
RUN go build -o tmp/main main.go


FROM golang:1.24-alpine
WORKDIR /app

RUN go install github.com/air-verse/air@latest
COPY --from=builder /app/tmp/main .
COPY --from=builder /app/.air.toml .
RUN mkdir -p /app/tmp /app/images && \
    chmod -R 777 /app/tmp /app/images
COPY --from=builder /app .

WORKDIR /app
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]