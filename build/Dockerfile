ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum first to leverage Docker cache
COPY ../go.mod ../go.sum ./
RUN go mod download && go mod verify

# Copy the entire project (all source files and assets)
COPY .. .

# Build the binary
RUN go build -v -o /lxkota cmd/server/main.go

FROM debian:bookworm

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /lxkota /app/

# Copy template file and other assets needed at runtime
COPY --from=builder /usr/src/app/web/template.html /app/web/template.html
COPY --from=builder /usr/src/app/web/static /app/web/static
COPY --from=builder /usr/src/app/web/static/ /app/static/
COPY --from=builder /usr/src/app/web/content /app/web/content

# Run the binary
CMD ["./lxkota"]