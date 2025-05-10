ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN ls -latr .
RUN go build -v -o /lxkota server/main.go

FROM debian:bookworm

COPY --from=builder /lxkota /usr/local/bin/
CMD ["lxkota"]
