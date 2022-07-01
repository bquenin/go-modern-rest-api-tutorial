FROM golang:1.19.0-alpine3.16 as builder
WORKDIR /work

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY api ./api
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o microservice ./cmd/microservice/main.go

FROM alpine:3.16.2
WORKDIR /bin
COPY --from=builder /work/microservice /bin/microservice
ENV GIN_MODE=release
CMD /bin/microservice
