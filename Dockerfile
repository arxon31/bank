FROM golang:1.21-alpine as builder
LABEL authors="arxon31"

WORKDIR /usr/local/src

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY . .
RUN go build -o ./bin/consumer ./cmd/consumer/main.go

# runner
FROM alpine AS runner

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/src/bin/consumer /usr/local/bin

CMD ["./consumer"]

