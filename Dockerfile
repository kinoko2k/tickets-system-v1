FROM golang:1.27-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY static ./static
COPY templates ./templates
COPY config.json .

RUN apk add --no-cache tzdata
ENV TZ=Asia/Tokyo

EXPOSE 3000
CMD ["./server"]
