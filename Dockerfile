FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN "go build -o /lxkota server/main.go

WORKDIR /

COPY --from=builder /lxkota /lxkota

USER nonroot:nonroot

ENTRYPOINT [ "/lxkota" ]