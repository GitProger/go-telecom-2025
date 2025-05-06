FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-telecom-2025 ./cmd/main.go

FROM alpine:latest AS runner
WORKDIR /root/
COPY --from=builder /app/go-telecom-2025 .
ENTRYPOINT ["./go-telecom-2025"]
# CMD [ "./sunny_5_skiers/config.json", "./sunny_5_skiers/events" ]
