# --- Stage 1: Build ---
FROM golang:1.25.5 AS builder

WORKDIR /app

COPY go.mod ./

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# --- Stage 2: Runtime ---
FROM scratch

WORKDIR /root/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/main .

CMD ["./main"]