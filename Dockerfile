FROM golang:1.22.2-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o server cmd/server/main.go

EXPOSE 8080

CMD ["sh", "./scripts/run_local.sh"]
