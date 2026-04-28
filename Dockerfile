FROM golang:1.25 AS builder

WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler
COPY --from=builder /app/web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["/app/scheduler"]