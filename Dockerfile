FROM golang:1.22-alpine3.18 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o confdi -ldflags="-s -w" main.go

FROM alpine:3.18 as runner

WORKDIR /app

COPY --from=builder /app/confdi .

ENTRYPOINT ["/app/confdi"]