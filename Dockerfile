FROM golang:1.24.0 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y git

COPY go.mod go.sum ./

RUN go mod tidy

COPY . ./

RUN ls -la

RUN GOOS=linux GOARCH=amd64 go build -o /app/bin/final_project ./backend/cmd/main.go

FROM ubuntu:latest


WORKDIR /app

RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/final_project .
COPY web /app/web
COPY scheduler.db .


ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db


EXPOSE ${TODO_PORT}

CMD ["./final_project"]