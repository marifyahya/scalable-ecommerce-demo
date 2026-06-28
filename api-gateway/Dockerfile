FROM golang:alpine

WORKDIR /app

# Install Air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Note: We don't COPY the rest of the files because docker-compose will mount the local directory as a volume

EXPOSE 8000

# Start Air
CMD ["air", "-c", ".air.toml"]
