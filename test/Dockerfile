FROM golang:1.19.0-alpine3.16
WORKDIR /root

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY api ./api
COPY cmd ./cmd
COPY internal ./internal

CMD CGO_ENABLED=0 go test ./...
